package entity

import "github.com/Mahoura-shop/Backend/internal/infrastructure/database"

type Province struct {
	database.Model
	Name   string `gorm:"type:varchar(50);not null"`
	Cities []City `gorm:"foreignKey:ProvinceID"`
}
