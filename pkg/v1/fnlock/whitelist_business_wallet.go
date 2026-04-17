package fnlock

// 钱包直付与 token 发送触发入口。
var whitelistBusinessWallet = []LockedFunction{
	{
		ID:               "bitfs.clientapp.wallet.trigger_pay_bsv",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerWalletBSVTransfer",
		Signature:        "func TriggerWalletBSVTransfer(ctx context.Context, store *clientDB, rt *Runtime, req WalletBSVTransferRequest) (WalletBSVTransferResult, error)",
		ObsControlAction: "wallet.pay_bsv",
		Note:             "钱包直连外部地址的 p2pkh 支付触发入口。",
	},
	{
		ID:               "bitfs.clientapp.wallet.trigger_token_send_preview",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerWalletTokenSendPreview",
		Signature:        "func TriggerWalletTokenSendPreview(ctx context.Context, store *clientDB, rt *Runtime, req WalletTokenSendPreviewRequest) (WalletAssetActionPreviewResp, error)",
		ObsControlAction: "wallet.token_send_preview",
		Note:             "bsv21 token 发送预演头部入口，供 HTTP 与 obs control 复用。",
	},
	{
		ID:               "bitfs.clientapp.wallet.trigger_token_send_sign",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerWalletTokenSendSign",
		Signature:        "func TriggerWalletTokenSendSign(ctx context.Context, store *clientDB, rt *Runtime, req WalletTokenSendSignRequest) (WalletAssetActionSignResp, error)",
		ObsControlAction: "wallet.token_send_sign",
		Note:             "bsv21 token 发送签名头部入口，强制校验 expected_preview_hash。",
	},
	{
		ID:               "bitfs.clientapp.wallet.trigger_token_send_submit",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerWalletTokenSendSubmit",
		Signature:        "func TriggerWalletTokenSendSubmit(ctx context.Context, store *clientDB, rt *Runtime, req WalletAssetActionSubmitRequest) (WalletAssetActionSubmitResp, error)",
		ObsControlAction: "wallet.token_send_submit",
		Note:             "bsv21 token 发送广播头部入口，统一本地投影与事实收口。",
	},
}
