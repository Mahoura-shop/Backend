package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type CartRepository struct{}

func NewCartRepository() *CartRepository {
	return &CartRepository{}
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

func (repo *CartRepository) CreateCart(db database.Database, cart *entity.Cart) error {
	return db.GetDB().Create(&cart).Error
}

func (repo *CartRepository) DeleteCartByID(db database.Database, id uint) error {
	return db.GetDB().Where("id = ?", id).Unscoped().Delete(&entity.Cart{}).Error
}

func (repo *CartRepository) UpdateCart(db database.Database, cart *entity.Cart) error {
	return db.GetDB().Save(&cart).Error
}