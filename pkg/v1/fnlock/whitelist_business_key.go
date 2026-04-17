package fnlock

// 密钥生命周期相关能力入口。
var whitelistBusinessKey = []LockedFunction{
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
}
