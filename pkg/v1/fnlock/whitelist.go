package fnlock

import (
	"sort"
	"strings"
)

// Module 表示函数归属模块。
type Module string

const (
	ModuleBitFS Module = "bitfs"
	ModuleBFTP  Module = "bftp"
)

// LockedFunction 是函数签名锁白名单项。
//
// 设计说明：
// - 进入白名单后，函数名、参数、返回值都被冻结；
// - CI 通过 go doc 提取实时签名，必须和 Signature 完全一致；
// - Note 必填，用于说明该签名为什么不能随意改动。
type LockedFunction struct {
	ID               string
	Module           Module
	Package          string
	Symbol           string
	Signature        string
	ObsControlAction string
	Note             string
}

// Whitelist 是跨 BitFS/BFTP 的全局函数签名锁清单。
var Whitelist = []LockedFunction{
	{
		ID:        "bitfs.managed.control.execute",
		Module:    ModuleBitFS,
		Package:   "./pkg/managedclient",
		Symbol:    "managedDaemon.executeManagedControlCommand",
		Signature: "func (d *managedDaemon) executeManagedControlCommand(req controlCommandRequest) controlCommandResult",
		Note:      "托管控制命令总分发口，e2e 与运行时都依赖该入口语义。",
	},
	{
		ID:        "bitfs.managed.control.execute_workspace",
		Module:    ModuleBitFS,
		Package:   "./pkg/managedclient",
		Symbol:    "managedDaemon.executeManagedWorkspaceControlCommand",
		Signature: "func (d *managedDaemon) executeManagedWorkspaceControlCommand(req controlCommandRequest) (controlCommandResult, error)",
		Note:      "workspace 控制动作唯一业务入口，避免中转壳与重复分发。",
	},
	{
		ID:        "bitfs.managed.control.execute_pricing",
		Module:    ModuleBitFS,
		Package:   "./pkg/managedclient",
		Symbol:    "managedDaemon.executeManagedPricingControlCommand",
		Signature: "func (d *managedDaemon) executeManagedPricingControlCommand(req controlCommandRequest) (controlCommandResult, error)",
		Note:      "定价控制桥接分发口，保持签名稳定可避免控制层与业务层脱钩。",
	},
	{
		ID:               "bitfs.managed.key.ensure_material",
		Module:           ModuleBitFS,
		Package:          "./pkg/managedclient",
		Symbol:           "managedDaemon.ensureKeyMaterial",
		Signature:        "func (d *managedDaemon) ensureKeyMaterial(password string) (keyMaterialTicket, error)",
		ObsControlAction: "key.ensure_material",
		Note:             "密钥材料创建链路的核心能力入口，测试触发与 HTTP 语义需要稳定对齐。",
	},
	{
		ID:               "bitfs.managed.key.unlock_with_password",
		Module:           ModuleBitFS,
		Package:          "./pkg/managedclient",
		Symbol:           "managedDaemon.unlockWithPassword",
		Signature:        "func (d *managedDaemon) unlockWithPassword(ctx context.Context, source, password string) (unlockTicket, error)",
		ObsControlAction: "key.unlock",
		Note:             "解锁链路总入口，锁状态/运行时启动/并发竞争都围绕该函数。",
	},
	{
		ID:               "bitfs.managed.key.lock_runtime",
		Module:           ModuleBitFS,
		Package:          "./pkg/managedclient",
		Symbol:           "managedDaemon.lockRuntime",
		Signature:        "func (d *managedDaemon) lockRuntime() error",
		ObsControlAction: "key.lock",
		Note:             "运行时落锁与停机收敛入口，签名变化会破坏状态机调用面。",
	},
	{
		ID:               "bitfs.clientapp.pricing.trigger_set_base",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerPricingSetBase",
		Signature:        "func TriggerPricingSetBase(ctx context.Context, rt *Runtime, base uint64) (PricingSetBaseResult, error)",
		ObsControlAction: "pricing.set_base",
		Note:             "定价基准价触发入口，控制层和业务层通过该触发器对接。",
	},
	{
		ID:               "bitfs.clientapp.pricing.trigger_reset_seed",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerPricingResetSeed",
		Signature:        "func TriggerPricingResetSeed(ctx context.Context, rt *Runtime, seedHash string) (PricingStateResult, error)",
		ObsControlAction: "pricing.reset_seed",
		Note:             "单种子定价状态重置触发入口。",
	},
	{
		ID:               "bitfs.clientapp.pricing.trigger_feed_seed",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerPricingFeedSeed",
		Signature:        "func TriggerPricingFeedSeed(ctx context.Context, rt *Runtime, req PricingFeedRequest) (PricingStateResult, error)",
		ObsControlAction: "pricing.feed_seed",
		Note:             "定价样本喂入触发入口。",
	},
	{
		ID:               "bitfs.clientapp.pricing.trigger_set_force",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerPricingSetForce",
		Signature:        "func TriggerPricingSetForce(ctx context.Context, rt *Runtime, req ForcePriceRequest) (PricingStateResult, error)",
		ObsControlAction: "pricing.set_force",
		Note:             "单种子强制定价触发入口。",
	},
	{
		ID:               "bitfs.clientapp.pricing.trigger_release_force",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerPricingReleaseForce",
		Signature:        "func TriggerPricingReleaseForce(ctx context.Context, rt *Runtime, seedHash string) (PricingStateResult, error)",
		ObsControlAction: "pricing.release_force",
		Note:             "释放强制定价触发入口。",
	},
	{
		ID:               "bitfs.clientapp.pricing.trigger_run_tick",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerPricingRunTick",
		Signature:        "func TriggerPricingRunTick(ctx context.Context, rt *Runtime, hours uint32) (PricingRunTickResult, error)",
		ObsControlAction: "pricing.run_tick",
		Note:             "定价 tick 执行触发入口。",
	},
	{
		ID:               "bitfs.clientapp.pricing.trigger_reconcile",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerPricingReconcile",
		Signature:        "func TriggerPricingReconcile(ctx context.Context, rt *Runtime, seedHash string, nowUnix int64) (PricingReconcileResult, error)",
		ObsControlAction: "pricing.trigger_reconcile",
		Note:             "单种子定价对账触发入口。",
	},
	{
		ID:               "bitfs.clientapp.pricing.trigger_get_config",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerPricingGetConfig",
		Signature:        "func TriggerPricingGetConfig(rt *Runtime) (PricingConfig, error)",
		ObsControlAction: "pricing.get_config",
		Note:             "读取定价全局配置触发入口。",
	},
	{
		ID:               "bitfs.clientapp.pricing.trigger_get_state",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerPricingGetState",
		Signature:        "func TriggerPricingGetState(ctx context.Context, rt *Runtime, seedHash string) (PricingState, error)",
		ObsControlAction: "pricing.get_state",
		Note:             "读取单种子定价状态触发入口。",
	},
	{
		ID:               "bitfs.clientapp.pricing.trigger_get_audits",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerPricingGetAudits",
		Signature:        "func TriggerPricingGetAudits(ctx context.Context, rt *Runtime, seedHash string, limit int) ([]PricingAuditItem, error)",
		ObsControlAction: "pricing.get_audits",
		Note:             "读取定价审计轨迹触发入口。",
	},
	{
		ID:               "bitfs.clientapp.pricing.trigger_list_seeds",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerPricingListSeeds",
		Signature:        "func TriggerPricingListSeeds(ctx context.Context, rt *Runtime, limit int) ([]string, error)",
		ObsControlAction: "pricing.list_seeds",
		Note:             "列举定价种子触发入口。",
	},
	{
		ID:               "bitfs.clientapp.workspace.kernel_list",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "workspaceManager.List",
		Signature:        "func (m *workspaceManager) List() ([]workspaceItem, error)",
		ObsControlAction: "workspace.list",
		Note:             "workspace 列表能力单点入口，调用侧统一从 workspace manager 获取结果。",
	},
	{
		ID:               "bitfs.clientapp.workspace.kernel_add",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "workspaceManager.Add",
		Signature:        "func (m *workspaceManager) Add(path string, maxBytes uint64) (workspaceItem, error)",
		ObsControlAction: "workspace.add",
		Note:             "workspace 新增能力单点入口。",
	},
	{
		ID:               "bitfs.clientapp.workspace.kernel_update",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "workspaceManager.UpdateByPath",
		Signature:        "func (m *workspaceManager) UpdateByPath(workspacePath string, maxBytes *uint64, enabled *bool) (workspaceItem, error)",
		ObsControlAction: "workspace.update",
		Note:             "workspace 更新能力单点入口。",
	},
	{
		ID:               "bitfs.clientapp.workspace.kernel_delete",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "workspaceManager.DeleteByPath",
		Signature:        "func (m *workspaceManager) DeleteByPath(workspacePath string) error",
		ObsControlAction: "workspace.delete",
		Note:             "workspace 删除能力单点入口。",
	},
	{
		ID:               "bitfs.clientapp.workspace.kernel_sync_once",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "workspaceManager.SyncOnce",
		Signature:        "func (m *workspaceManager) SyncOnce(ctx context.Context) (map[string]sellerSeed, error)",
		ObsControlAction: "workspace.sync_once",
		Note:             "workspace 单次扫描同步能力单点入口。",
	},
	{
		ID:        "bftp.obs.info",
		Module:    ModuleBFTP,
		Package:   "./pkg/obs",
		Symbol:    "Info",
		Signature: "func Info(service, name string, fields map[string]any)",
		Note:      "系统流程日志统一入口。",
	},
	{
		ID:        "bftp.obs.debug",
		Module:    ModuleBFTP,
		Package:   "./pkg/obs",
		Symbol:    "Debug",
		Signature: "func Debug(service, name string, fields map[string]any)",
		Note:      "调试日志统一入口。",
	},
	{
		ID:        "bftp.obs.business",
		Module:    ModuleBFTP,
		Package:   "./pkg/obs",
		Symbol:    "Business",
		Signature: "func Business(service, name string, fields map[string]any)",
		Note:      "业务日志统一入口。",
	},
	{
		ID:        "bftp.obs.important",
		Module:    ModuleBFTP,
		Package:   "./pkg/obs",
		Symbol:    "Important",
		Signature: "func Important(service, name string, fields map[string]any)",
		Note:      "关键业务日志统一入口。",
	},
	{
		ID:        "bftp.obs.error",
		Module:    ModuleBFTP,
		Package:   "./pkg/obs",
		Symbol:    "Error",
		Signature: "func Error(service, name string, fields map[string]any)",
		Note:      "错误日志统一入口。",
	},
	{
		ID:        "bftp.sqliteactor.open",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/sqliteactor",
		Symbol:    "Open",
		Signature: "func Open(path string, debug bool) (*Opened, error)",
		Note:      "sqlite 单 owner actor 的标准打开入口。",
	},
	{
		ID:        "bftp.sqliteactor.new",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/sqliteactor",
		Symbol:    "New",
		Signature: "func New(db *sql.DB) (*Actor, error)",
		Note:      "actor 构造入口，串行访问边界依赖该签名。",
	},
	{
		ID:        "bftp.sqliteactor.actor_do",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/sqliteactor",
		Symbol:    "Actor.Do",
		Signature: "func (a *Actor) Do(ctx context.Context, fn func(*sql.DB) error) error",
		Note:      "sqlite 非事务执行主入口。",
	},
	{
		ID:        "bftp.sqliteactor.actor_tx",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/sqliteactor",
		Symbol:    "Actor.Tx",
		Signature: "func (a *Actor) Tx(ctx context.Context, fn func(*sql.Tx) error) error",
		Note:      "sqlite 事务执行主入口。",
	},
	{
		ID:        "bftp.sqliteactor.do_value",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/sqliteactor",
		Symbol:    "DoValue",
		Signature: "func DoValue[T any](ctx context.Context, a *Actor, fn func(*sql.DB) (T, error)) (T, error)",
		Note:      "sqlite 非事务泛型返回值入口。",
	},
	{
		ID:        "bftp.sqliteactor.tx_value",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/sqliteactor",
		Symbol:    "TxValue",
		Signature: "func TxValue[T any](ctx context.Context, a *Actor, fn func(*sql.Tx) (T, error)) (T, error)",
		Note:      "sqlite 事务泛型返回值入口。",
	},
	{
		ID:        "bftp.ncall.register",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/ncall",
		Symbol:    "Register",
		Signature: "func Register(h host.Host, sec pproto.SecurityConfig, callHandler CallHandler, resolveHandler ResolveHandler)",
		Note:      "统一 node.call / node.resolve 协议注册入口。",
	},
}

var obsControlActionToLockID = buildObsControlActionToLockID()

func buildObsControlActionToLockID() map[string]string {
	out := map[string]string{}
	for _, item := range Whitelist {
		action := strings.TrimSpace(item.ObsControlAction)
		lockID := strings.TrimSpace(item.ID)
		if action == "" || lockID == "" {
			continue
		}
		out[action] = lockID
	}
	return out
}

// ObsControlActions 返回 obs 控制动作列表（已排序）。
func ObsControlActions() []string {
	out := make([]string, 0, len(obsControlActionToLockID))
	for action := range obsControlActionToLockID {
		out = append(out, action)
	}
	sort.Strings(out)
	return out
}

// IsObsControlActionAllowed 判断 obs 控制动作是否在白名单内。
func IsObsControlActionAllowed(action string) bool {
	_, ok := ObsControlActionLockID(action)
	return ok
}

// ObsControlActionLockID 返回动作对应的函数签名锁 ID。
func ObsControlActionLockID(action string) (string, bool) {
	lockID, ok := obsControlActionToLockID[strings.TrimSpace(action)]
	return lockID, ok
}
