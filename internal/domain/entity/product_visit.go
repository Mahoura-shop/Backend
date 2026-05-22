package entity

import (
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type ProductVisit struct {
	database.Model
	ProductID uint   `gorm:"index"`
	Product   Product `gorm:"foreignKey:ProductID"`
	VisitorIP string `gorm:"type:varchar(45);index"`
	VisitedAt int64 `gorm:"index"`
}
