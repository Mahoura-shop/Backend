package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Order struct {
	database.Model
	CartID    uint
	Cart      Cart    `gorm:"foreignKey:CartID"`
	PaymentID uint
	Payment   Payment `gorm:"foreignKey:PaymentID"`
}