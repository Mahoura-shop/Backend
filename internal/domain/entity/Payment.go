package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Payment struct {
	database.Model
	Amount uint
}