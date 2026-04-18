package fnlock

// peer.* 触发入口。
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
	{
		ID:               "bitfs.clientapp.peer.trigger_self",
		Module:           ModuleBitFS,
		Package:          "./pkg/managedclient",
		Symbol:           "managedDaemon.executeManagedBusinessControlCommand",
		Signature:        "func (d *managedDaemon) executeManagedBusinessControlCommand(req controlCommandRequest) (controlCommandResult, error)",
		ObsControlAction: "peer.self",
		Note:             "managed 控制入口：读取当前节点 peer 身份（pubkey_hex/peer_id/addrs）。",
	},
	{
		ID:               "bitfs.clientapp.peer.trigger_connect",
		Module:           ModuleBitFS,
		Package:          "./pkg/managedclient",
		Symbol:           "managedDaemon.executeManagedBusinessControlCommand",
		Signature:        "func (d *managedDaemon) executeManagedBusinessControlCommand(req controlCommandRequest) (controlCommandResult, error)",
		ObsControlAction: "peer.connect",
		Note:             "managed 控制入口：按给定 addr 建立 p2p 连接（写 peerstore + host connect）。",
	},
}
