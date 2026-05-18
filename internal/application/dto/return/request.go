package returndto

type RequestReturnRequest struct {
	UserID      uint   `json:"-"`
	OrderItemID uint   `json:"orderItemID" binding:"required"`
	Reason      string `json:"reason" binding:"required"`
	Quantity    uint   `json:"quantity" binding:"required,min=1"`
}

type ReviewReturnRequest struct {
	AdminID uint   `json:"-"`
	Action  string `json:"action" binding:"required,oneof=approve reject"`
	Note    string `json:"note"`
}
