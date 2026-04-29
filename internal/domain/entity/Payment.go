package entity

import (
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Payment struct {
	database.Model
	OrderID    uint
	Order      Order              `gorm:"foreignKey:OrderID"`
	Amount     uint               `gorm:"type:int"`
	Status     enum.PaymentStatus `gorm:"index"`
	Authority  string             `gorm:"type:varchar(100);uniqueIndex"`
	RefCode    string             `gorm:"type:varchar(100)"`
	GatewayURL string             `gorm:"type:varchar(255)"`
}
