package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type OrderRepository struct{}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{}
}

func (repo *OrderRepository) FindOrderByID(db database.Database, orderID uint) (*entity.Order, error) {
	var order entity.Order
	result := db.GetDB().Preload("Brand").Preload("Category").Preload("Currency").Where("id = ?", orderID).First(&order)
	if result.Error != nil {
		return nil, result.Error
	}
	return &order, nil
}

func (repo *OrderRepository) GetOrders(db database.Database) ([]*entity.Order, error) {
	var orders []*entity.Order
	result := db.GetDB().Preload("Brand").Preload("Category").Preload("Currency").Find(&orders)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return orders, nil
}

func (repo *OrderRepository) CreateOrder(db database.Database) (*entity.Order, error) {
	order := entity.Order{}
    result := db.GetDB().Create(&order)
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &order, nil
}

func (repo *OrderRepository) DeleteOrderByID(db database.Database, orderID uint) (error) {
	return db.GetDB().Where("id = ?", orderID).Unscoped().Delete(&entity.Order{}).Error
}



func (repo *OrderRepository) CreateOrderItem(db database.Database, orderItem entity.OrderItem) (error) {
    result := db.GetDB().Create(&orderItem)
    if result.Error != nil {
        return result.Error
    }
    
    return nil
}