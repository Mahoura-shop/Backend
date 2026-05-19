package entity

import "github.com/Mahoura-shop/Backend/internal/infrastructure/database"

type AdminActivityLog struct {
	database.Model
	AdminID    uint   `gorm:"not null;index"`
	Admin      User   `gorm:"foreignKey:AdminID"`
	Method     string `gorm:"type:varchar(10);not null"`
	Path       string `gorm:"type:varchar(500);not null"`
	IPAddress  string `gorm:"type:varchar(50)"`
	StatusCode int
}
