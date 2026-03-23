package entity

import (
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Cart struct {
	database.Model
	UserID uint
	User   User            `gorm:"foreignKey:UserID"`
	Status enum.CartStatus `gorm:"index"`
	Items  []Product       `gorm:"many2many:cart_items;constraint:OnDelete:CASCADE;"`
}