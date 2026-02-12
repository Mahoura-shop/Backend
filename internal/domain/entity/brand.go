package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Brand struct {
	database.Model
	Name        string `gorm:"type:varchar(50);not null"`
	Slug        string `gorm:"type:varchar(50);uniqueIndex"`
	Description string `gorm:"type:text"`
	IsActive    bool   `gorm:"default:true"`
	BrandPic    string `gorm:"type:varchar(255);default:null"`
}
