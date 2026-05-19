package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type RoleRepository struct{}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{}
}

func (r *RoleRepository) GetRoles(db database.Database) ([]*entity.Role, error) {
	var roles []*entity.Role
	if err := db.GetDB().Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

func (r *RoleRepository) FindRoleByID(db database.Database, id uint) (*entity.Role, error) {
	var role entity.Role
	if err := db.GetDB().Preload("Permissions").First(&role, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) CreateRole(db database.Database, role entity.Role) (*entity.Role, error) {
	if err := db.GetDB().Create(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) UpdateRole(db database.Database, role entity.Role) error {
	perms := role.Permissions
	role.Permissions = nil
	if err := db.GetDB().Save(&role).Error; err != nil {
		return err
	}
	return db.GetDB().Model(&role).Association("Permissions").Replace(perms)
}

func (r *RoleRepository) DeleteRoleByID(db database.Database, id uint) error {
	var role entity.Role
	if err := db.GetDB().First(&role, id).Error; err != nil {
		return err
	}
	if err := db.GetDB().Model(&entity.User{}).Where("role_id = ?", id).Update("role_id", nil).Error; err != nil {
		return err
	}
	if err := db.GetDB().Model(&role).Association("Permissions").Clear(); err != nil {
		return err
	}
	return db.GetDB().Unscoped().Delete(&role).Error
}

func (r *RoleRepository) GetPermissions(db database.Database) ([]*entity.Permission, error) {
	var permissions []*entity.Permission
	if err := db.GetDB().Find(&permissions).Error; err != nil {
		return nil, err
	}
	return permissions, nil
}
