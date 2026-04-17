package fnlock

// domain 交互触发入口。
var whitelistBusinessDomain = []LockedFunction{
	{
		ID:               "bitfs.clientapp.domain.trigger_resolve",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerResolverResolve",
		Signature:        "func TriggerResolverResolve(ctx context.Context, store *clientDB, rt *Runtime, p TriggerResolverResolveParams) (resolverResolveResp, error)",
		ObsControlAction: "domain.resolve",
		Note:             "domain name 到目标公钥解析触发入口。",
	},
	{
		ID:               "bitfs.clientapp.domain.trigger_register",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerDomainRegisterName",
		Signature:        "func TriggerDomainRegisterName(ctx context.Context, store *clientDB, rt *Runtime, p TriggerDomainRegisterNameParams) (TriggerDomainRegisterNameResult, error)",
		ObsControlAction: "domain.register",
		Note:             "domain 注册全流程触发入口（query -> lock -> submit）。",
	},
	{
		ID:               "bitfs.clientapp.domain.trigger_set_target",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerDomainSetTarget",
		Signature:        "func TriggerDomainSetTarget(ctx context.Context, store *clientDB, rt *Runtime, p TriggerDomainSetTargetParams) (TriggerDomainSetTargetResult, error)",
		ObsControlAction: "domain.set_target",
		Note:             "domain 修改 target 指向触发入口。",
	},
}
