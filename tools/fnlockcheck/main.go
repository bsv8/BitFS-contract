package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bsv8/bitfs-contract/pkg/v1/fnlock"
)

// moduleConfig 描述一个可检查模块在工作区中的落盘目录。
// 这里用最小信息建映射：模块名 -> 子目录名。
type moduleConfig struct {
	name string
	dir  string
}

// moduleDirs 把 fnlock 中的模块常量映射到真实源码目录。
// runChecks 会按这个映射切到对应目录执行 go doc。
var moduleDirs = map[fnlock.Module]moduleConfig{
	fnlock.ModuleBitFS: {name: string(fnlock.ModuleBitFS), dir: "BitFS"},
	fnlock.ModuleBFTP:  {name: string(fnlock.ModuleBFTP), dir: "BFTP"},
}

// main 是 fnlock 校验工具入口，执行顺序固定：
// 1) 解析参数与工作区
// 2) 解析目标模块
// 3) 校验白名单结构是否合法
// 4) 校验 obs 控制动作绑定是否合法
// 5) 对每条白名单做签名比对
// 任一步失败都立即退出并打印统一前缀错误。
func main() {
	var (
		workspaceRootFlag string
		goBin             string
		modulesFlag       string
	)
	flag.StringVar(&workspaceRootFlag, "workspace-root", "", "workspace root path (contains go.work)")
	flag.StringVar(&goBin, "go-bin", "/home/david/.gvm/gos/go1.26.0/bin/go", "go binary path")
	flag.StringVar(&modulesFlag, "modules", "bitfs,bftp", "comma-separated modules to check")
	flag.Parse()

	workspaceRoot, err := resolveWorkspaceRoot(workspaceRootFlag)
	if err != nil {
		exitErr(err)
	}
	selectedModules, err := parseModules(modulesFlag)
	if err != nil {
		exitErr(err)
	}
	if err := validateWhitelistShape(fnlock.Whitelist); err != nil {
		exitErr(err)
	}
	if err := validateObsControlBindings(fnlock.Whitelist); err != nil {
		exitErr(err)
	}
	if err := runChecks(workspaceRoot, goBin, selectedModules, fnlock.Whitelist); err != nil {
		exitErr(err)
	}
	fmt.Println("[fnlock] ok")
}

// resolveWorkspaceRoot 解析工作区根目录（要求目录内有 go.work）。
// 优先使用传入参数；如果没传，就从当前目录开始向上逐级查找。
func resolveWorkspaceRoot(input string) (string, error) {
	if v := strings.TrimSpace(input); v != "" {
		abs, err := filepath.Abs(v)
		if err != nil {
			return "", err
		}
		if _, err := os.Stat(filepath.Join(abs, "go.work")); err != nil {
			return "", fmt.Errorf("workspace root missing go.work: %s", abs)
		}
		return abs, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(wd, "go.work")); err == nil {
			return wd, nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	return "", errors.New("cannot find workspace root with go.work")
}

// parseModules 解析 -modules 参数，转成集合形式，方便后续 O(1) 判断是否需要检查。
// 例如 "bitfs,bftp" -> {bitfs:{}, bftp:{}}。
func parseModules(raw string) (map[fnlock.Module]struct{}, error) {
	out := map[fnlock.Module]struct{}{}
	for _, item := range strings.Split(raw, ",") {
		m := fnlock.Module(strings.TrimSpace(strings.ToLower(item)))
		if m == "" {
			continue
		}
		if _, ok := moduleDirs[m]; !ok {
			return nil, fmt.Errorf("unsupported module: %s", item)
		}
		out[m] = struct{}{}
	}
	if len(out) == 0 {
		return nil, errors.New("no module selected")
	}
	return out, nil
}

// validateWhitelistShape 只做“白名单自身形状”校验，不碰源码：
// - 必填字段是否为空
// - id 是否重复
// - module 是否支持
// - signature 是否以 func 开头（保证是函数签名）
// - module/package/symbol 组合是否重复（防止一处函数被重复定义规则）
func validateWhitelistShape(items []fnlock.LockedFunction) error {
	seenID := map[string]struct{}{}
	seenSymbol := map[string]struct{}{}
	for i, item := range items {
		prefix := fmt.Sprintf("whitelist[%d]", i)
		if strings.TrimSpace(item.ID) == "" {
			return fmt.Errorf("%s id is required", prefix)
		}
		if _, ok := seenID[item.ID]; ok {
			return fmt.Errorf("%s duplicated id: %s", prefix, item.ID)
		}
		seenID[item.ID] = struct{}{}
		if _, ok := moduleDirs[item.Module]; !ok {
			return fmt.Errorf("%s unsupported module: %s", prefix, item.Module)
		}
		if strings.TrimSpace(item.Package) == "" {
			return fmt.Errorf("%s package is required", prefix)
		}
		if strings.TrimSpace(item.Symbol) == "" {
			return fmt.Errorf("%s symbol is required", prefix)
		}
		if strings.TrimSpace(item.Signature) == "" {
			return fmt.Errorf("%s signature is required", prefix)
		}
		if !strings.HasPrefix(strings.TrimSpace(item.Signature), "func ") {
			return fmt.Errorf("%s signature must start with func", prefix)
		}
		if strings.TrimSpace(item.Note) == "" {
			return fmt.Errorf("%s note is required", prefix)
		}
		symbolKey := string(item.Module) + "|" + item.Package + "|" + item.Symbol
		if _, ok := seenSymbol[symbolKey]; ok {
			return fmt.Errorf("%s duplicated module/package/symbol: %s", prefix, symbolKey)
		}
		seenSymbol[symbolKey] = struct{}{}
	}
	return nil
}

// validateObsControlBindings 校验“obs 控制动作 <-> 锁条目”关系是否合法：
// - 一个 action 只能出现一次，避免路由歧义
// - 只要设置了 ObsControlAction，就必须有 lock id
// - 当前策略要求 action 仅绑定 BitFS 侧函数（避免跨模块误绑定）
func validateObsControlBindings(items []fnlock.LockedFunction) error {
	seenAction := map[string]struct{}{}
	for i, item := range items {
		prefix := fmt.Sprintf("whitelist[%d]", i)
		action := strings.TrimSpace(item.ObsControlAction)
		if action == "" {
			continue
		}
		if _, ok := seenAction[action]; ok {
			return fmt.Errorf("%s duplicated action: %s", prefix, action)
		}
		seenAction[action] = struct{}{}
		lockID := strings.TrimSpace(item.ID)
		if lockID == "" {
			return fmt.Errorf("%s lock_id is required when obs_control_action is set", prefix)
		}
		if item.Module != fnlock.ModuleBitFS {
			return fmt.Errorf("%s lock_id must point to bitfs function: %s", prefix, lockID)
		}
	}
	return nil
}

// runChecks 执行真正的签名冻结校验：
// - 遍历白名单
// - 对选中模块条目调用 go doc 读取“当前真实签名”
// - 与白名单中的期望签名逐条做字符串严格比较
// 最终把所有失败项聚合后一次性返回，便于一次看全改动影响面。
func runChecks(workspaceRoot string, goBin string, selected map[fnlock.Module]struct{}, items []fnlock.LockedFunction) error {
	goBin = strings.TrimSpace(goBin)
	if goBin == "" {
		return errors.New("go-bin is required")
	}
	goBinAbs, err := filepath.Abs(goBin)
	if err != nil {
		return err
	}
	if _, err := os.Stat(goBinAbs); err != nil {
		return fmt.Errorf("go binary not found: %s", goBinAbs)
	}
	goRoot := filepath.Dir(filepath.Dir(goBinAbs))
	pathPrefix := filepath.Dir(goBinAbs)

	var failed []string
	for _, item := range items {
		// 只校验用户选中的模块，支持按模块分批检查。
		if _, ok := selected[item.Module]; !ok {
			continue
		}
		cfg := moduleDirs[item.Module]
		moduleDir := filepath.Join(workspaceRoot, cfg.dir)
		got, err := readSignature(goBinAbs, goRoot, pathPrefix, moduleDir, item.Package, item.Symbol)
		if err != nil {
			failed = append(failed, fmt.Sprintf("%s: %v", item.ID, err))
			continue
		}
		if strings.TrimSpace(got) != strings.TrimSpace(item.Signature) {
			failed = append(failed, fmt.Sprintf("%s: signature mismatch\n  want: %s\n  got:  %s", item.ID, item.Signature, got))
		}
	}
	if len(failed) == 0 {
		return nil
	}
	sort.Strings(failed)
	return fmt.Errorf("function lock check failed:\n- %s", strings.Join(failed, "\n- "))
}

// readSignature 在指定模块目录下执行：
//   go doc -u <pkg> <symbol>
// 然后从输出中提取第一条以 "func " 开头的行，作为真实签名。
// 这里不解析 AST，直接以 go doc 结果作为“外部可见声明”的事实来源。
func readSignature(goBin string, goRoot string, goBinDir string, moduleDir string, pkg string, symbol string) (string, error) {
	cmd := exec.Command(goBin, "doc", "-u", strings.TrimSpace(pkg), strings.TrimSpace(symbol))
	cmd.Dir = moduleDir
	cmd.Env = enrichEnv(goRoot, goBinDir)
	out, err := cmd.CombinedOutput()
	text := strings.TrimSpace(string(out))
	if err != nil {
		if text == "" {
			text = err.Error()
		}
		return "", fmt.Errorf("go doc failed: %s", text)
	}
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "func ") {
			return line, nil
		}
	}
	return "", fmt.Errorf("cannot find func signature from go doc output: %s", text)
}

// enrichEnv 为 go doc 命令补齐可预期运行环境：
// - PATH 前置 goBinDir，保证优先使用传入 go 工具链
// - 可选设置 GOROOT，避免机器上多版本 go 时解析漂移
func enrichEnv(goRoot string, goBinDir string) []string {
	base := os.Environ()
	pathValue := os.Getenv("PATH")
	if strings.TrimSpace(pathValue) == "" {
		pathValue = goBinDir
	} else {
		pathValue = goBinDir + string(os.PathListSeparator) + pathValue
	}
	base = append(base, "PATH="+pathValue)
	if strings.TrimSpace(goRoot) != "" {
		base = append(base, "GOROOT="+goRoot)
	}
	return base
}

// exitErr 统一错误出口，保证脚本侧拿到稳定前缀日志。
func exitErr(err error) {
	fmt.Fprintf(os.Stderr, "[fnlock] %s\n", strings.TrimSpace(err.Error()))
	os.Exit(1)
}
