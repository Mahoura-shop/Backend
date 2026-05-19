package entity

import (
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Notification struct {
	database.Model
	UserID uint                 `gorm:"not null;index"`
	Type   enum.NotificationType `gorm:"type:smallint;not null"`
	Title  string               `gorm:"not null"`
	Body   string               `gorm:"not null"`
	IsRead bool                 `gorm:"default:false"`
	RefID  *uint
}
