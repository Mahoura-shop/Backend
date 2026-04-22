package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Currency struct {
	database.Model
	Name        string `gorm:"type:varchar(50);uniqueIndex"`
	Code        string `gorm:"type:varchar(5);default:'IRR';uniqueIndex"`
	ConvertRate uint   `gorm:"type:int;not null"`
}