package entity

import "github.com/Mahoura-shop/Backend/internal/infrastructure/database"

type Province struct {
	database.Model
	Name         string `gorm:"type:varchar(50);not null"`
	ShippingCost uint   `gorm:"type:int;default:0"`
	Cities       []City `gorm:"foreignKey:ProvinceID"`
}
