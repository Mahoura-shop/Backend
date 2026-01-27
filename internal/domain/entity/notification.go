package entity

import "github.com/CosmeticsShiraz/Backend/internal/infrastructure/database"

type Notification struct {
	database.Model
	TypeID      uint             `gorm:"not null;index"`
	Type        NotificationType `gorm:"foreignKey:TypeID"`
	IsRead      bool             `gorm:"default:false"`
	RecipientID uint             `gorm:"not null;index"`
	Data        []byte           `gorm:"type:jsonb"`
}
