package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type RevenueByTier struct {
	Tier    string
	Revenue uint
}

type OrderByDay struct {
	Date  string
	Count uint
}

type OrderRepository interface {
	FindOrderByID(db database.Database, orderID uint) (*entity.Order, error)
	GetOrders(db database.Database) ([]*entity.Order, error)
	GetOrdersByUserID(db database.Database, userID uint) ([]*entity.Order, error)
	CreateOrder(db database.Database, order entity.Order) (*entity.Order, error)
	UpdateOrder(db database.Database, order entity.Order) error
	DeleteOrderByID(db database.Database, orderID uint) error
	CreateOrderItem(db database.Database, orderItem entity.OrderItem) error
	CreateOrderStatusHistory(db database.Database, history entity.OrderStatusHistory) error
	GetRevenuePerTier(db database.Database) ([]RevenueByTier, error)
	GetOrdersPerDay(db database.Database, days int) ([]OrderByDay, error)
}
