package entity

import (
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type User struct {
	database.Model
	FirstName      string          `gorm:"type:varchar(50);index:idx_user_name"`
	LastName       string          `gorm:"type:varchar(50);index:idx_user_name"`
	Phone          string          `gorm:"type:varchar(20);uniqueIndex"`
	Email          string          `gorm:"type:varchar(100);Index"`
	EmailVerified  bool            `gorm:"default:false"`
	ProfilePicPath string          `gorm:"type:varchar(255);default:null"`
	Status         enum.UserStatus `gorm:"index"`
	Addresses      []Address       `gorm:"polymorphic:Owner;polymorphicValue:users"`
	IsAdmin        bool            `gorm:"default:false"`
	Type           enum.UserType   `gorm:"index"`
}
