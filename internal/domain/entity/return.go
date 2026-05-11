package entity

import (
	"time"

	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Return struct {
	database.Model
	OrderItemID uint             `gorm:"not null;index"`
	OrderItem   OrderItem        `gorm:"foreignKey:OrderItemID"`
	UserID      uint             `gorm:"not null;index"`
	User        User             `gorm:"foreignKey:UserID"`
	Status      enum.ReturnStatus `gorm:"index"`
	Reason      string           `gorm:"type:text"`
	Quantity    uint             `gorm:"not null"`
	RefundAmount uint            `gorm:"not null"`
	RequestedAt time.Time        `gorm:"autoCreateTime"`
	ApprovedAt  *time.Time
	RefundedAt  *time.Time
	ApprovedByID *uint
	ApprovedBy  *User `gorm:"foreignKey:ApprovedByID"`
}
