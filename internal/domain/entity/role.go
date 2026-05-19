package entity

import "github.com/Mahoura-shop/Backend/internal/infrastructure/database"

type Permission struct {
	database.Model
	Name        string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Description string `gorm:"type:varchar(255)"`
	Category    string `gorm:"type:varchar(50);not null;default:''"`
}

type Role struct {
	database.Model
	Name        string       `gorm:"type:varchar(100);uniqueIndex;not null"`
	Description string       `gorm:"type:varchar(255)"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
