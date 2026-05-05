package entity

import "github.com/Mahoura-shop/Backend/internal/infrastructure/database"

type ProductImage struct {
	database.Model
	ProductID uint   `gorm:"not null;index"`
	Path      string `gorm:"type:varchar(255);not null"`
	Position  uint   `gorm:"default:0"`
}
