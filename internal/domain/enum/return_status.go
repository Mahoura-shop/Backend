package enum

type ReturnStatus string

const (
	ReturnStatusRequested  ReturnStatus = "requested"
	ReturnStatusApproved   ReturnStatus = "approved"
	ReturnStatusRejected   ReturnStatus = "rejected"
	ReturnStatusShippedBack ReturnStatus = "shipped_back"
	ReturnStatusReceived   ReturnStatus = "received"
	ReturnStatusRefunded   ReturnStatus = "refunded"
)
