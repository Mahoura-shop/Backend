package entity

import (
	"time"

	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Instalment struct {
	database.Model
	OrderID  uint
	Order    Order                 `gorm:"foreignKey:OrderID"`
	Number   uint                 `gorm:"type:int;not null"`
	Amount   uint                 `gorm:"type:int;not null"`
	DueDate  time.Time            `gorm:"not null"`
	Status   enum.InstalmentStatus `gorm:"type:smallint;default:1"`
	PaidAt   *time.Time
}
