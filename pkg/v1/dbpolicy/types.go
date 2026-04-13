package dbpolicy

// SettlementFactTable 是链上支付事实表名。
// 设计约束：
// - quote/pay 链路可以写入；
// - 其他业务链路只读，不能直接写事实表。
const (
	TableFactSettlementChannelChainQuotePay    = "fact_settlement_channel_chain_quote_pay"
	TableFactSettlementChannelChainDirectPay   = "fact_settlement_channel_chain_direct_pay"
	TableFactSettlementChannelChainAssetCreate = "fact_settlement_channel_chain_asset_create"
)

// WritableByQuotePay 明确 quote/pay 链路允许写入的事实表。
var WritableByQuotePay = map[string]struct{}{
	TableFactSettlementChannelChainQuotePay: {},
}

// CanWriteByQuotePay 判断 quote/pay 链路是否允许写目标事实表。
func CanWriteByQuotePay(table string) bool {
	_, ok := WritableByQuotePay[table]
	return ok
}
