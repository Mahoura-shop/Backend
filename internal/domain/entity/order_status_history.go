package entity

import (
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type OrderStatusHistory struct {
	database.Model
	OrderID     uint
	Order       Order            `gorm:"foreignKey:OrderID"`
	Status      enum.OrderStatus `gorm:"type:smallint"`
	Note        string           `gorm:"type:text"`
	ChangedByID uint
}
