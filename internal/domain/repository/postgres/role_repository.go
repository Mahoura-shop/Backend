package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type RoleRepository interface {
	GetRoles(database.Database) ([]*entity.Role, error)
	FindRoleByID(database.Database, uint) (*entity.Role, error)
	CreateRole(database.Database, entity.Role) (*entity.Role, error)
	UpdateRole(database.Database, entity.Role) error
	DeleteRoleByID(database.Database, uint) error
	GetPermissions(database.Database) ([]*entity.Permission, error)
}
