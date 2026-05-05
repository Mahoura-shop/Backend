package entity

import "github.com/Mahoura-shop/Backend/internal/infrastructure/database"

type Review struct {
	database.Model
	UserID     uint    `gorm:"uniqueIndex:idx_review_user_product;not null"`
	User       User    `gorm:"foreignKey:UserID"`
	ProductID  uint    `gorm:"uniqueIndex:idx_review_user_product;not null"`
	Rating     uint    `gorm:"type:smallint;not null;check:rating >= 1 AND rating <= 5"`
	Comment    string  `gorm:"type:text"`
	IsVerified bool    `gorm:"default:false"`
}
