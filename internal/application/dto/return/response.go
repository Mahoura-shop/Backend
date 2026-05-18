package returndto

import "time"

type ReturnCredential struct {
	ID           uint       `json:"id"`
	OrderItemID  uint       `json:"orderItemID"`
	ProductName  string     `json:"productName"`
	UserID       uint       `json:"userID"`
	UserPhone    string     `json:"userPhone,omitempty"`
	Status       string     `json:"status"`
	Reason       string     `json:"reason"`
	Quantity     uint       `json:"quantity"`
	RefundAmount uint       `json:"refundAmount"`
	RequestedAt  time.Time  `json:"requestedAt"`
	ApprovedAt   *time.Time `json:"approvedAt,omitempty"`
	RefundedAt   *time.Time `json:"refundedAt,omitempty"`
}
