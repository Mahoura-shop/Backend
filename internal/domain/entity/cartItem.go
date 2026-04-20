package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type CartItem struct {
    database.Model
    CartID    uint
    Cart      Cart    `gorm:"foreignKey:CartID"`
    ProductID uint
    Product   Product `gorm:"foreignKey:ProductID"`
    Count     uint    `gorm:"type:int"`
}