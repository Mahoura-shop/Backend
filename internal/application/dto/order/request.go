package orderdto

import "github.com/Mahoura-shop/Backend/internal/domain/enum"

type CreateOrderRequest struct {
	UserID              uint
	AddressID           *uint              `json:"addressID"`
	PaymentMethod       enum.PaymentMethod `json:"paymentMethod"`
	InstalmentCount     uint               `json:"instalmentCount"`
	InstalmentIntervalDays uint            `json:"instalmentIntervalDays"`
}

type UpdateOrderStatusRequest struct {
	Status      enum.OrderStatus `json:"status" validate:"required"`
	Note        string           `json:"note"`
	ChangedByID uint
}
