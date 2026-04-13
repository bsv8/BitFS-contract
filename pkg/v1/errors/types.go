package bitfserrors

// Code 是 bitfs-contract 层统一错误码。
type Code string

const (
	CodeOK               Code = "OK"
	CodeBadRequest       Code = "BAD_REQUEST"
	CodeUnauthorized     Code = "UNAUTHORIZED"
	CodeNotFound         Code = "NOT_FOUND"
	CodeConflict         Code = "CONFLICT"
	CodeTimeout          Code = "TIMEOUT"
	CodeRetryable        Code = "RETRYABLE"
	CodeNonRetryable     Code = "NON_RETRYABLE"
	CodeChainTxFailed    Code = "CHAIN_TX_FAILED"
	CodePaymentRequired  Code = "PAYMENT_REQUIRED"
	CodeRouteNotFound    Code = "ROUTE_NOT_FOUND"
	CodeInternal         Code = "INTERNAL"
	CodeUnspecifiedError Code = "UNSPECIFIED_ERROR"
)

// Fault 是契约层错误对象。
// 说明：Message 保持英文，便于跨语言一致处理。
type Fault struct {
	Code    Code
	Message string
}

func (f Fault) Error() string {
	if f.Message == "" {
		return string(f.Code)
	}
	return string(f.Code) + ": " + f.Message
}
