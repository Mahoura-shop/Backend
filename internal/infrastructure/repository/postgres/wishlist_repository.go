package postgres

import (
	"errors"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type WishlistRepository struct{}

func NewWishlistRepository() *WishlistRepository {
	return &WishlistRepository{}
}

func (r *WishlistRepository) AddToWishlist(db database.Database, item entity.Wishlist) error {
	return db.GetDB().Create(&item).Error
}

func (r *WishlistRepository) RemoveFromWishlist(db database.Database, userID, productID uint) error {
	return db.GetDB().Unscoped().Where("user_id = ? AND product_id = ?", userID, productID).Delete(&entity.Wishlist{}).Error
}

func (r *WishlistRepository) GetWishlistByUserID(db database.Database, userID uint) ([]*entity.Wishlist, error) {
	var items []*entity.Wishlist
	result := db.GetDB().
		Preload("Product").
		Preload("Product.Brand").
		Preload("Product.Category").
		Preload("Product.Currency").
		Where("user_id = ?", userID).
		Find(&items)
	if result.Error != nil {
		return nil, result.Error
	}
	return items, nil
}

func (r *WishlistRepository) FindWishlistItem(db database.Database, userID, productID uint) (*entity.Wishlist, error) {
	var item entity.Wishlist
	result := db.GetDB().Where("user_id = ? AND product_id = ?", userID, productID).First(&item)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &item, nil
}
