package entity

import (
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Payment struct {
	database.Model
	Amount uint               `gorm:"type:int"`
	Status enum.PaymentStatus `gorm:"index"`
}