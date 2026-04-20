package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Order struct {
	database.Model
	Items []OrderItem `gorm:"constraint:OnDelete:CASCADE;"`
}