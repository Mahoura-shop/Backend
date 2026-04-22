package usecase

import (
	orderdto "github.com/Mahoura-shop/Backend/internal/application/dto/order"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
)

type OrderService interface {
	ParseOrder(entity.Order) (orderdto.OrderCredential)
	GetOrder(uint) (*orderdto.OrderCredential, error)
	GetOrders() ([]orderdto.OrderCredential, error)
	RegisterOrder(uint) (error)
}
