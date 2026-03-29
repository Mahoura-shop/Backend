package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type OrderItem struct {
	database.Model
    ProductID uint
    Product       Product `gorm:"foreignKey:ProductID"`
    Count         uint    `gorm:"type:int"`
	PriceSnapshot uint    `gorm:"type:int"`
}