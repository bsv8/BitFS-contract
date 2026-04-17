package fnlock

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
