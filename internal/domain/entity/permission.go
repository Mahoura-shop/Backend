package entity

import (
	"github.com/CosmeticsShiraz/Backend/internal/domain/enum"
	"github.com/CosmeticsShiraz/Backend/internal/infrastructure/database"
)

type Permission struct {
	database.Model
	Type        enum.PermissionType `gorm:"not null;index"`
	Description string              `gorm:"type:text"`
	Category    enum.PermissionCategory
	Roles       []Role `gorm:"many2many:role_permissions;constraint:OnDelete:CASCADE;"`
}
