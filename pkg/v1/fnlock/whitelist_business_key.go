package fnlock

// 密钥生命周期相关能力入口。
var whitelistBusinessKey = []LockedFunction{
	{
		ID:               "bitfs.managed.key.create_random",
		Module:           ModuleBitFS,
		Package:          "./pkg/managedclient",
		Symbol:           "managedDaemon.createRandomKeyMaterial",
		Signature:        "func (d *managedDaemon) createRandomKeyMaterial(password string) (keyMaterialTicket, error)",
		ObsControlAction: "key.create_random",
		Note:             "只负责创建随机密钥；已存在时直接报错，避免测试语义混淆。",
	},
	{
		ID:               "bitfs.managed.key.assert_exists",
		Module:           ModuleBitFS,
		Package:          "./pkg/managedclient",
		Symbol:           "managedDaemon.assertKeyMaterialExists",
		Signature:        "func (d *managedDaemon) assertKeyMaterialExists() error",
		ObsControlAction: "key.assert_exists",
		Note:             "只负责断言密钥文件存在；不创建、不修改。",
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
