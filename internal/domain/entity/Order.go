package entity

import (
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Order struct {
	database.Model
	UserID        uint
	User          User               `gorm:"foreignKey:UserID"`
	PaymentMethod enum.PaymentMethod `gorm:"type:smallint"`
	Items         []OrderItem        `gorm:"constraint:OnDelete:CASCADE;"`
}