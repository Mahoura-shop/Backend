package entity

import (
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type Permission struct {
	database.Model
	Type        enum.PermissionType `gorm:"not null;index"`
	Description string              `gorm:"type:text"`
	Category    enum.PermissionCategory
	UserType    enum.UserType
	Roles       []Role `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE;"`
}
