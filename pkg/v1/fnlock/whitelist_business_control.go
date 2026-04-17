package fnlock

// 托管控制分发入口。
var whitelistBusinessControl = []LockedFunction{
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
}
