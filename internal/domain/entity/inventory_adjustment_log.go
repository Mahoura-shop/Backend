package entity

import (
	"time"

	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type InventoryAdjustmentLog struct {
	database.Model
	ProductID uint      `gorm:"not null;index"`
	Product   Product   `gorm:"foreignKey:ProductID"`
	UserID    uint      `gorm:"not null;index"`
	User      User      `gorm:"foreignKey:UserID"`
	OldQuantity uint    `gorm:"not null"`
	NewQuantity uint    `gorm:"not null"`
	Reason     string   `gorm:"type:text"`
	AdjustedAt time.Time `gorm:"autoCreateTime"`
}
