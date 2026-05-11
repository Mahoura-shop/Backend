package usecase

import roledto "github.com/Mahoura-shop/Backend/internal/application/dto/role"

type RoleService interface {
	GetRoles() ([]roledto.RoleCredential, error)
	GetPermissions() ([]roledto.PermissionCredential, error)
	CreateRole(roledto.CreateRoleRequest) (*roledto.RoleCredential, error)
	UpdateRole(roledto.UpdateRoleRequest) error
	DeleteRole(id uint) error
}
