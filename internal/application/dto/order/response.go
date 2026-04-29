package orderdto

import (
	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
)

type OrderCredential struct {
	ID            uint                  `json:"id"`
	UserID        uint                  `json:"userID"`
	PaymentMethod enum.PaymentMethod    `json:"paymentMethod"`
	Items         []OrderItemCredential `json:"items"`
}

type OrderItemCredential struct {
	ID            uint                         `json:"id"`
	Product       productdto.ProductCredential `json:"product"`
	Count         uint                         `json:"count"`
	PriceSnapshot uint                         `json:"priceSnapshot"`
	Tier          enum.UserType                `json:"tier"`
}
