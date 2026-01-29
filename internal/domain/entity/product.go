package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Product struct {
	database.Model
	Name         string    `gorm:"type:varchar(50);not null"`
	Slug         string    `gorm:"type:varchar(50);uniqueIndex"`
	Description  string   `gorm:"type:text"`
	IsActive     bool      `gorm:"default:true"`
	IsNew        bool      `gorm:"default:true"`
	Priority     uint      `gorm:"default:0;index" validate:"min=0"`
	MinOrder     uint      `gorm:"default:1" validate:"min=1"`
	CategoryID   *uint
	Category     *Category `gorm:"foreignKey:CategoryID"`
	Quantity     uint      `gorm:"default:0;not null"`
	QuantityType string   `gorm:"default:'pieces';not null"`
	Price        float64   `gorm:"type:decimal(10,2);not null"`
	CurrencyCode string    `gorm:"type:varchar(5);default:'IRR';not null"`
}
