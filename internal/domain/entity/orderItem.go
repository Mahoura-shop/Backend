package entity

import (
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type OrderItem struct {
	database.Model
	OrderID       uint
	Order         Order         `gorm:"foreignKey:OrderID"`
	ProductID     uint
	Product       Product       `gorm:"foreignKey:ProductID"`
	Count         uint          `gorm:"type:int"`
	PriceSnapshot uint          `gorm:"type:int;not null"`
	Tier          enum.UserType `gorm:"type:smallint;not null"`
}