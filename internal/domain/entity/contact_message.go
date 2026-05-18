package entity

import "github.com/Mahoura-shop/Backend/internal/infrastructure/database"

type ContactMessage struct {
	database.Model
	Name    string `gorm:"not null"`
	Email   string `gorm:"not null;index"`
	Subject string `gorm:"not null"`
	Message string `gorm:"type:text;not null"`
}
