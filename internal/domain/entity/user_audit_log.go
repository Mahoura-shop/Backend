package entity

import (
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type UserAuditLog struct {
	database.Model
	UserID      uint
	User        User          `gorm:"foreignKey:UserID"`
	ChangedByID uint
	OldType     enum.UserType `gorm:"type:smallint"`
	NewType     enum.UserType `gorm:"type:smallint"`
	Reason      string        `gorm:"type:text"`
}
