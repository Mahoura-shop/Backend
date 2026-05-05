package entity

import "github.com/Mahoura-shop/Backend/internal/infrastructure/database"

type Wishlist struct {
	database.Model
	UserID    uint    `gorm:"uniqueIndex:idx_wishlist_user_product;not null"`
	ProductID uint    `gorm:"uniqueIndex:idx_wishlist_user_product;not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
}
