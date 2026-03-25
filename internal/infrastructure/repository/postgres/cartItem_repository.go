package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type CartItemRepository struct{}

func NewCartItemRepository() *CartItemRepository {
	return &CartItemRepository{}
} 

func (repo *CartItemRepository) FindCartItemByID(db database.Database, cartItemID uint) (*entity.CartItem, error) {
	var cartItem entity.CartItem
	result := db.GetDB().First(&cartItem, cartItemID)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &cartItem, nil
}

func (repo *CartItemRepository) FindCartItemByProductID(db database.Database, productID uint) (*entity.CartItem, error) {
	var cartItem entity.CartItem
	result := db.GetDB().Where("product_id = ?", productID).First(&cartItem)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &cartItem, nil
}

func (repo *CartItemRepository) CreateCartItem(db database.Database, cartItem entity.CartItem) (error) {
	return db.GetDB().Create(&cartItem).Error
}

func (repo *CartItemRepository) DeleteCartItemByID(db database.Database, cartItemID uint) (error) {
	return db.GetDB().Where("id = ?", cartItemID).Unscoped().Delete(&entity.CartItem{}).Error
}

func (repo *CartItemRepository) UpdateCartItem(db database.Database, cartItem entity.CartItem) (error) {
	return db.GetDB().Save(&cartItem).Error
}

func (repo *CartItemRepository) IncreaseCartItemCount(db database.Database, cartItemID uint) (error) {
	var cartItem entity.CartItem
	result := db.GetDB().First(&cartItem, cartItemID)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}
	}
	cartItem.Count++;
	err := db.GetDB().Save(&cartItem).Error
	return err
}

func (repo *CartItemRepository) DecreaseCartItemCount(db database.Database, cartItemID uint) (error) {
	var cartItem entity.CartItem
	result := db.GetDB().First(&cartItem, cartItemID)
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}
	}
	cartItem.Count--;
	err := db.GetDB().Save(&cartItem).Error
	return err
}