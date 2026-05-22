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
        Preload("Items").
        Preload("Items.Product").
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

func (repo *CartRepository) FindCartByID(db database.Database, cartID uint) (*entity.Cart, error) {
	var cart entity.Cart
	result := db.GetDB().First(&cart, cartID)
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

func (repo *CartRepository) DeleteCart(db database.Database, cartID uint) (error) {
	return db.GetDB().Where("id = ?", cartID).Delete(&entity.Cart{}).Error
}

func (repo *CartRepository) DeleteCartItems(db database.Database, cartID uint) (error) {
    return db.GetDB().Where("cart_id = ?", cartID).Delete(&entity.CartItem{}).Error;
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
        return nil
    }

    return err
}

func (repo *CartRepository) FindCartItemByID(db database.Database, cartItemID uint) (*entity.CartItem, error) {
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

func (repo *CartRepository) FindCartItemByProductID(db database.Database, productID uint, cartID uint) (*entity.CartItem, error) {
	var cartItem entity.CartItem
	result := db.GetDB().Where("product_id = ? AND cart_id = ?", productID, cartID).First(&cartItem)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &cartItem, nil
}

func (repo *CartRepository) CreateCartItem(db database.Database, cartItem entity.CartItem) (error) {
	return db.GetDB().Create(&cartItem).Error
}

func (repo *CartRepository) DeleteCartItemByID(db database.Database, cartItemID uint) (error) {
	return db.GetDB().Where("id = ?", cartItemID).Unscoped().Delete(&entity.CartItem{}).Error
}

func (repo *CartRepository) UpdateCartItem(db database.Database, cartItem entity.CartItem) (error) {
	return db.GetDB().Save(&cartItem).Error
}

func (repo *CartRepository) IncreaseCartItemCount(db database.Database, cartItemID uint) (error) {
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

func (repo *CartRepository) DecreaseCartItemCount(db database.Database, cartItemID uint) (error) {
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

func (repo *CartRepository) DeleteCartItemsByProductID(db database.Database, productID uint) (error) {
	return db.GetDB().Where("product_id = ?", productID).Unscoped().Delete(&entity.CartItem{}).Error
}