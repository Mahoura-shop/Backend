package postgres

import (
	"errors"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
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
