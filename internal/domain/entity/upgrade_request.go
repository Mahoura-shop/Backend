package entity

import (
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type UpgradeRequest struct {
	database.Model
	UserID        uint
	User          User                      `gorm:"foreignKey:UserID"`
	RequestedType enum.UserType             `gorm:"type:smallint;not null"`
	BusinessName  string                    `gorm:"type:varchar(200)"`
	TaxID         string                    `gorm:"type:varchar(50)"`
	StorePhoto    string                    `gorm:"type:varchar(255)"`
	Status        enum.UpgradeRequestStatus `gorm:"type:smallint;not null;default:1"`
	AdminNote     string                    `gorm:"type:text"`
	ReviewedByID  *uint
	ReviewedBy    *User                     `gorm:"foreignKey:ReviewedByID"`
}
