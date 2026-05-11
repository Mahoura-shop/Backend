package entity

import (
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Ticket struct {
	database.Model
	UserID    uint               `gorm:"not null;index"`
	User      User               `gorm:"foreignKey:UserID"`
	Subject   enum.TicketSubject `gorm:"index"`
	Status    enum.TicketStatus  `gorm:"index"`
	Message   string             `gorm:"type:text"`
	Response  string             `gorm:"type:text"`
	OrderID   *uint
	Order     *Order `gorm:"foreignKey:OrderID"`
}
