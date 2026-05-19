package postgres

import (
	"errors"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	postgresrepo "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (repo *OrderRepository) FindOrderByID(db database.Database, orderID uint) (*entity.Order, error) {
	var order entity.Order
	result := db.GetDB().
		Preload("User").
		Preload("Address").
		Preload("Items").
		Preload("Items.Product").
		Preload("StatusHistory").
		Where("id = ?", orderID).
		First(&order)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (repo *OrderRepository) GetOrders(db database.Database) ([]*entity.Order, error) {
	var orders []*entity.Order
	result := db.GetDB().
		Preload("User").
		Preload("Address").
		Preload("Items").
		Preload("Items.Product").
		Order("created_at DESC").
		Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func (repo *OrderRepository) GetOrdersByUserID(db database.Database, userID uint) ([]*entity.Order, error) {
	var orders []*entity.Order
	result := db.GetDB().
		Preload("Address").
		Preload("Items").
		Preload("Items.Product").
		Preload("StatusHistory").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders)
	if result.Error != nil {
		return nil, result.Error
	}
	return orders, nil
}

func (repo *OrderRepository) CreateOrder(db database.Database, order entity.Order) (*entity.Order, error) {
	result := db.GetDB().Create(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (repo *OrderRepository) UpdateOrder(db database.Database, order entity.Order) error {
	return db.GetDB().Save(&order).Error
}

func (repo *OrderRepository) DeleteOrderByID(db database.Database, orderID uint) error {
	return db.GetDB().Where("id = ?", orderID).Unscoped().Delete(&entity.Order{}).Error
}

func (repo *OrderRepository) CreateOrderItem(db database.Database, orderItem entity.OrderItem) error {
	return db.GetDB().Create(&orderItem).Error
}

func (repo *OrderRepository) CreateOrderStatusHistory(db database.Database, history entity.OrderStatusHistory) error {
	return db.GetDB().Create(&history).Error
}

func (repo *OrderRepository) GetRevenuePerTier(db database.Database) ([]postgresrepo.RevenueByTier, error) {
	var results []postgresrepo.RevenueByTier
	result := db.GetDB().
		Model(&entity.OrderItem{}).
		Select(`
			CASE tier
				WHEN 2 THEN 'regular'
				WHEN 3 THEN 'shopkeeper'
				WHEN 4 THEN 'shopkeeper'
				WHEN 5 THEN 'fellow'
				WHEN 6 THEN 'admin'
			END as tier,
			SUM(price_snapshot * count) as revenue
		`).
		Group("tier").
		Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}
	return results, nil
}

func (repo *OrderRepository) GetRevenuePerDay(db database.Database, days int) ([]postgresrepo.RevenueByDay, error) {
	var results []postgresrepo.RevenueByDay
	result := db.GetDB().
		Model(&entity.OrderItem{}).
		Select("DATE(orders.created_at) as date, SUM(order_items.price_snapshot * order_items.count) as revenue").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.created_at >= NOW() - ? * INTERVAL '1 day'", days).
		Group("DATE(orders.created_at)").
		Order("DATE(orders.created_at)").
		Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}
	return results, nil
}

func (repo *OrderRepository) FindOrderItemByID(db database.Database, orderItemID uint) (*entity.OrderItem, error) {
	var item entity.OrderItem
	result := db.GetDB().Preload("Order").Preload("Product").Where("id = ?", orderItemID).First(&item)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, result.Error
}

func (repo *OrderRepository) GetOrdersPerDay(db database.Database, days int) ([]postgresrepo.OrderByDay, error) {
	var results []postgresrepo.OrderByDay
	result := db.GetDB().
		Model(&entity.Order{}).
		Select("DATE(created_at) as date, COUNT(*) as count").
		Where("created_at >= NOW() - ? * INTERVAL '1 day'", days).
		Group("DATE(created_at)").
		Order("DATE(created_at)").
		Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}
	return results, nil
}
