package entity

import (
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Order struct {
	database.Model
	UserID        uint
	User          User                 `gorm:"foreignKey:UserID"`
	AddressID     *uint
	Address       *Address             `gorm:"foreignKey:AddressID"`
	PaymentMethod enum.PaymentMethod   `gorm:"type:smallint"`
	Status        enum.OrderStatus     `gorm:"type:smallint;default:1"`
	TotalAmount   uint                 `gorm:"type:int"`
	ShippingCost  uint                 `gorm:"type:int"`
	RefundFlag    bool                 `gorm:"default:false"`
	Items         []OrderItem          `gorm:"constraint:OnDelete:CASCADE;"`
	StatusHistory []OrderStatusHistory `gorm:"foreignKey:OrderID"`
}
