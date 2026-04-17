package fnlock

// 节点网络协议入口。
var whitelistBusinessNetwork = []LockedFunction{
	{
		ID:        "bftp.ncall.register",
		Module:    ModuleBFTP,
		Package:   "./pkg/infra/ncall",
		Symbol:    "Register",
		Signature: "func Register(h host.Host, sec pproto.SecurityConfig, callHandler CallHandler, resolveHandler ResolveHandler)",
		Note:      "统一 node.call / node.resolve 协议注册入口。",
	},
}
