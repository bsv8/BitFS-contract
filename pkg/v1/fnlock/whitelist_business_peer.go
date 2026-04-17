package fnlock

// peer.call / peer.resolve 触发入口。
var whitelistBusinessPeer = []LockedFunction{
	{
		ID:               "bitfs.clientapp.peer.trigger_call",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerPeerCall",
		Signature:        "func TriggerPeerCall(ctx context.Context, rt *Runtime, p TriggerPeerCallParams) (ncall.CallResp, error)",
		ObsControlAction: "peer.call",
		Note:             "节点 route call 触发入口，统一承接自动支付/quote/pay 模式。",
	},
	{
		ID:               "bitfs.clientapp.peer.trigger_resolve",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerPeerResolve",
		Signature:        "func TriggerPeerResolve(ctx context.Context, rt *Runtime, p TriggerPeerResolveParams) (ncall.ResolveResp, error)",
		ObsControlAction: "peer.resolve",
		Note:             "节点 route resolve 触发入口。",
	},
}
