package orderdto

import "github.com/Mahoura-shop/Backend/internal/domain/enum"

type CreateOrderRequest struct {
	UserID        uint
	PaymentMethod enum.PaymentMethod `json:"paymentMethod"`
}

type CreateOrderItemRequest struct {
	ProductID     uint
	Count         uint
	PriceSnapshot uint
}
