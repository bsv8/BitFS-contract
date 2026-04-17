package fnlock

// 网关交互触发入口。
var whitelistBusinessGateway = []LockedFunction{
	{
		ID:               "bitfs.clientapp.gateway.trigger_publish_demand",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerGatewayPublishDemand",
		Signature:        "func TriggerGatewayPublishDemand(ctx context.Context, store *clientDB, rt *Runtime, p PublishDemandParams) (contractmessage.DemandPublishPaidResp, error)",
		ObsControlAction: "gateway.publish_demand",
		Note:             "网关 demand 单条发布触发入口。",
	},
	{
		ID:               "bitfs.clientapp.gateway.trigger_publish_demand_batch",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerGatewayPublishDemandBatch",
		Signature:        "func TriggerGatewayPublishDemandBatch(ctx context.Context, store *clientDB, rt *Runtime, p PublishDemandBatchParams) (contractmessage.DemandPublishBatchPaidResp, error)",
		ObsControlAction: "gateway.publish_demand_batch",
		Note:             "网关 demand 批量发布触发入口。",
	},
	{
		ID:               "bitfs.clientapp.gateway.trigger_publish_live_demand",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerGatewayPublishLiveDemand",
		Signature:        "func TriggerGatewayPublishLiveDemand(ctx context.Context, store *clientDB, rt *Runtime, p PublishLiveDemandParams) (contractmessage.LiveDemandPublishPaidResp, error)",
		ObsControlAction: "gateway.publish_live_demand",
		Note:             "网关 live demand 发布触发入口。",
	},
	{
		ID:               "bitfs.clientapp.gateway.trigger_publish_demand_chain_tx_quote_pay",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerGatewayDemandPublishChainTxQuotePay",
		Signature:        "func TriggerGatewayDemandPublishChainTxQuotePay(ctx context.Context, store ClientStore, env gatewayDemandPublishChainTxEnv, p PublishDemandParams) (TriggerGatewayDemandPublishChainTxQuotePayResult, error)",
		ObsControlAction: "gateway.publish_demand_chain_tx_quote_pay",
		Note:             "网关 demand 的 chain_tx_v1 直付触发入口。",
	},
	{
		ID:               "bitfs.clientapp.gateway.trigger_reachability_announce",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerGatewayAnnounceNodeReachability",
		Signature:        "func TriggerGatewayAnnounceNodeReachability(ctx context.Context, store *clientDB, rt *Runtime, p AnnounceNodeReachabilityParams) (contractmessage.NodeReachabilityAnnouncePaidResp, error)",
		ObsControlAction: "gateway.reachability_announce",
		Note:             "节点向网关发布可达地址声明触发入口。",
	},
	{
		ID:               "bitfs.clientapp.gateway.trigger_reachability_query",
		Module:           ModuleBitFS,
		Package:          "./pkg/clientapp",
		Symbol:           "TriggerGatewayQueryNodeReachability",
		Signature:        "func TriggerGatewayQueryNodeReachability(ctx context.Context, store *clientDB, rt *Runtime, p QueryNodeReachabilityParams) (contractmessage.NodeReachabilityQueryPaidResp, error)",
		ObsControlAction: "gateway.reachability_query",
		Note:             "向网关查询节点可达地址声明触发入口。",
	},
}
