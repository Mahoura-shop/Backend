package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Wallet struct {
	database.Model
	Balance uint `gorm:"type:int"`
	UserID  uint
	User    User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}