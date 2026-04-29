package orderdto

import (
	"time"

	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
)

type OrderCredential struct {
	ID            uint                      `json:"id"`
	UserID        uint                      `json:"userID"`
	AddressID     *uint                     `json:"addressID"`
	PaymentMethod enum.PaymentMethod        `json:"paymentMethod"`
	Status        enum.OrderStatus          `json:"status"`
	TotalAmount   uint                      `json:"totalAmount"`
	ShippingCost  uint                      `json:"shippingCost"`
	RefundFlag    bool                      `json:"refundFlag"`
	Items         []OrderItemCredential     `json:"items"`
	StatusHistory []OrderStatusHistoryEntry `json:"statusHistory,omitempty"`
	CreatedAt     time.Time                 `json:"createdAt"`
}

type OrderItemCredential struct {
	ID            uint                         `json:"id"`
	Product       productdto.ProductCredential `json:"product"`
	Count         uint                         `json:"count"`
	PriceSnapshot uint                         `json:"priceSnapshot"`
	Tier          enum.UserType                `json:"tier"`
}

type OrderStatusHistoryEntry struct {
	Status      enum.OrderStatus `json:"status"`
	Note        string           `json:"note"`
	ChangedByID uint             `json:"changedByID"`
	CreatedAt   time.Time        `json:"createdAt"`
}

type InstalmentCredential struct {
	ID      uint                  `json:"id"`
	OrderID uint                  `json:"orderID"`
	Number  uint                  `json:"number"`
	Amount  uint                  `json:"amount"`
	DueDate time.Time             `json:"dueDate"`
	Status  enum.InstalmentStatus `json:"status"`
	PaidAt  *time.Time            `json:"paidAt"`
}

type PaymentGatewayResponse struct {
	GatewayURL string `json:"gatewayURL"`
}
