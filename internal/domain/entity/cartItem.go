package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type CartItem struct {
    database.Model
    CartID    uint
    ProductID uint
    Product   Product `gorm:"foreignKey:ProductID"`
    Count     uint    `gorm:"type:int"`
}