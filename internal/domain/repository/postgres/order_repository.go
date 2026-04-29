package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type OrderRepository interface {
	FindOrderByID(db database.Database, orderID uint) (*entity.Order, error)
	GetOrders(db database.Database) ([]*entity.Order, error)
	CreateOrder(db database.Database, order entity.Order) (*entity.Order, error)
	DeleteOrderByID(db database.Database, orderID uint) (error)

	CreateOrderItem(db database.Database, orderItem entity.OrderItem) (error)
}
