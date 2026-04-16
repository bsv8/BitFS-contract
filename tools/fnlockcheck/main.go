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

type moduleConfig struct {
	name string
	dir  string
}

var moduleDirs = map[fnlock.Module]moduleConfig{
	fnlock.ModuleBitFS: {name: string(fnlock.ModuleBitFS), dir: "BitFS"},
	fnlock.ModuleBFTP:  {name: string(fnlock.ModuleBFTP), dir: "BFTP"},
}

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

func exitErr(err error) {
	fmt.Fprintf(os.Stderr, "[fnlock] %s\n", strings.TrimSpace(err.Error()))
	os.Exit(1)
}
