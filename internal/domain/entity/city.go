package entity

import "github.com/CosmeticsShiraz/Backend/internal/infrastructure/database"

type City struct {
	database.Model
	Name       string `gorm:"type:varchar(50);not null"`
	ProvinceID uint   `gorm:"not null"`
}
