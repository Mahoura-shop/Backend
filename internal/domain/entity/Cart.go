package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Cart struct {
	database.Model
	UserID uint
	User   User       `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Items  []CartItem `gorm:"constraint:OnDelete:CASCADE;"`
}