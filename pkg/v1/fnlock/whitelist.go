package fnlock

// Whitelist 是跨 BitFS/BFTP 的全局函数签名锁清单。
//
// 设计说明：
// - 这里只做业务组装，不承载具体业务项；
// - 新增业务时，优先新增对应业务分组文件，再在这里挂接。
var Whitelist = buildWhitelist()

func buildWhitelist() []LockedFunction {
	total := len(whitelistBusinessControl) +
		len(whitelistBusinessKey) +
		len(whitelistBusinessWallet) +
		len(whitelistBusinessGateway) +
		len(whitelistBusinessDomain) +
		len(whitelistBusinessPeer) +
		len(whitelistBusinessPricing) +
		len(whitelistBusinessWorkspace) +
		len(whitelistBusinessObs) +
		len(whitelistBusinessSQLiteActor) +
		len(whitelistBusinessNetwork)
	out := make([]LockedFunction, 0, total)
	out = append(out, whitelistBusinessControl...)
	out = append(out, whitelistBusinessKey...)
	out = append(out, whitelistBusinessWallet...)
	out = append(out, whitelistBusinessGateway...)
	out = append(out, whitelistBusinessDomain...)
	out = append(out, whitelistBusinessPeer...)
	out = append(out, whitelistBusinessPricing...)
	out = append(out, whitelistBusinessWorkspace...)
	out = append(out, whitelistBusinessObs...)
	out = append(out, whitelistBusinessSQLiteActor...)
	out = append(out, whitelistBusinessNetwork...)
	return out
}
