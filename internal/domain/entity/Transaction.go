package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Transaction struct {
	database.Model
	Amount   int    `gorm:"type:int"`
	WalletID uint 
	Wallet   Wallet `gorm:"foreignKey:WalletID"`
}