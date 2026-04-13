package protoid

import "github.com/libp2p/go-libp2p/core/protocol"

// BitFS v0 首批协议 ID（直接传输 + 直播）。
// 设计说明：
// - 这里只放跨项目共享的协议常量；
// - 运行时代码统一引用这里，避免协议字面量散落在实现层。
const (
	ProtoHealth    protocol.ID = "/bsv-transfer/healthz/1.0.0"
	ProtoArbHealth protocol.ID = "/bsv-transfer/arbiter/healthz/1.0.0"

	ProtoSeedGet           protocol.ID = "/bsv-transfer/client/seed/get/1.0.0"
	ProtoQuoteDirectSubmit protocol.ID = "/bsv-transfer/client/quote/direct_submit/1.0.0"
	ProtoDirectDealAccept  protocol.ID = "/bsv-transfer/client/deal/accept/1.0.0"
	ProtoTransferPoolOpen  protocol.ID = "/bsv-transfer/client/transfer-pool/open/1.0.0"
	ProtoTransferPoolPay   protocol.ID = "/bsv-transfer/client/transfer-pool/pay/1.0.0"
	ProtoTransferPoolClose protocol.ID = "/bsv-transfer/client/transfer-pool/close/1.0.0"

	ProtoLiveQuoteSubmit protocol.ID = "/bsv-transfer/client/live_quote/submit/1.0.0"
	ProtoLiveSubscribe   protocol.ID = "/bsv-transfer/live/subscribe/1.0.0"
	ProtoLiveHeadPush    protocol.ID = "/bsv-transfer/live/head-push/1.0.0"
)
