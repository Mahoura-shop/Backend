package postgres

import (
	"errors"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type CartRepository struct{}

func NewCartRepository() *CartRepository {
	return &CartRepository{}
}

func (repo *CartRepository) FindCartByUserID(db database.Database, userID uint) (*entity.Cart, error) {
    var cart entity.Cart
    result := db.GetDB().
        Preload("User").
        Where("user_id = ?", userID).
        First(&cart)
    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, nil
    }
    if result.Error != nil {
        return nil, result.Error
    }
    return &cart, nil
}

func (repo *CartRepository) FindCartByID(db database.Database, id uint) (*entity.Cart, error) {
	var cart entity.Cart
	result := db.GetDB().First(&cart, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &cart, nil
}

func (repo *CartRepository) CreateCart(db database.Database, cart entity.Cart) (error) {
	return db.GetDB().Create(&cart).Error
}

func (repo *CartRepository) DeleteCartByID(db database.Database, id uint) (error) {
	return db.GetDB().Where("id = ?", id).Unscoped().Delete(&entity.Cart{}).Error
}

func (repo *CartRepository) UpdateCart(db database.Database, cart entity.Cart) (error) {
	return db.GetDB().Save(&cart).Error
}

func (repo *CartRepository) AddProductToCart(db database.Database, productID uint, cartID uint) (error) {
    var item entity.CartItem
    err := db.GetDB().
        Where("cart_id = ? AND product_id = ?", cartID, productID).
        First(&item).Error

    if err == nil {
        item.Count++
        return db.GetDB().Save(&item).Error
    }

    if errors.Is(err, gorm.ErrRecordNotFound) {
        item = entity.CartItem{
            CartID:    cartID,
            ProductID: productID,
            Count:     1,
        }
        return db.GetDB().Create(&item).Error
    }

    return err
}

func (repo *CartRepository) RemoveProductToCart(db database.Database, productID uint, cartID uint) (error) {
    var item entity.CartItem
    err := db.GetDB().
        Where("cart_id = ? AND product_id = ?", cartID, productID).
        First(&item).Error

    if err == nil {
        item.Count--
        return db.GetDB().Save(&item).Error
    }

    if errors.Is(err, gorm.ErrRecordNotFound) {
        item = entity.CartItem{
            CartID:    cartID,
            ProductID: productID,
            Count:     1,
        }
        return db.GetDB().Create(&item).Error
    }

    return err
}