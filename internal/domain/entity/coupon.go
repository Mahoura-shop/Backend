package entity

import (
	"time"

	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Coupon struct {
	database.Model
	Code            string     `gorm:"type:varchar(50);uniqueIndex;not null"`
	DiscountPercent float64    `gorm:"type:decimal(5,2);not null"`
	MinIRRPrice     uint       `gorm:"type:int;default:0"`
	MaxUses         uint       `gorm:"type:int;default:0"`
	UsedCount       uint       `gorm:"type:int;default:0"`
	ExpiresAt       *time.Time
	IsActive        bool       `gorm:"default:true"`
}
