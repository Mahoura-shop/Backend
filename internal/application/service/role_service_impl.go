package service

import (
	roledto "github.com/Mahoura-shop/Backend/internal/application/dto/role"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type RoleService struct {
	roleRepository postgres.RoleRepository
	db             database.Database
}

type RoleServiceDeps struct {
	RoleRepository postgres.RoleRepository
	DB             database.Database
}

func NewRoleService(deps RoleServiceDeps) *RoleService {
	return &RoleService{
		roleRepository: deps.RoleRepository,
		db:             deps.DB,
	}
}

func (s *RoleService) GetRoles() ([]roledto.RoleCredential, error) {
	roles, err := s.roleRepository.GetRoles(s.db)
	if err != nil {
		return nil, err
	}
	result := make([]roledto.RoleCredential, len(roles))
	for i, r := range roles {
		result[i] = mapRole(r)
	}
	return result, nil
}

func (s *RoleService) GetPermissions() ([]roledto.PermissionCredential, error) {
	perms, err := s.roleRepository.GetPermissions(s.db)
	if err != nil {
		return nil, err
	}
	result := make([]roledto.PermissionCredential, len(perms))
	for i, p := range perms {
		result[i] = roledto.PermissionCredential{ID: p.ID, Name: p.Name, Description: p.Description}
	}
	return result, nil
}

func (s *RoleService) CreateRole(req roledto.CreateRoleRequest) (*roledto.RoleCredential, error) {
	perms := make([]entity.Permission, len(req.PermissionIDs))
	for i, id := range req.PermissionIDs {
		perms[i] = entity.Permission{}
		perms[i].ID = id
	}
	role := entity.Role{
		Name:        req.Name,
		Description: req.Description,
		Permissions: perms,
	}
	created, err := s.roleRepository.CreateRole(s.db, role)
	if err != nil {
		return nil, err
	}
	cred := mapRole(created)
	return &cred, nil
}

func (s *RoleService) UpdateRole(req roledto.UpdateRoleRequest) error {
	existing, err := s.roleRepository.FindRoleByID(s.db, req.ID)
	if err != nil || existing == nil {
		return err
	}
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	perms := make([]entity.Permission, len(req.PermissionIDs))
	for i, id := range req.PermissionIDs {
		perms[i] = entity.Permission{}
		perms[i].ID = id
	}
	existing.Permissions = perms
	return s.roleRepository.UpdateRole(s.db, *existing)
}

func (s *RoleService) DeleteRole(id uint) error {
	return s.roleRepository.DeleteRoleByID(s.db, id)
}

func mapRole(r *entity.Role) roledto.RoleCredential {
	perms := make([]roledto.PermissionCredential, len(r.Permissions))
	for i, p := range r.Permissions {
		perms[i] = roledto.PermissionCredential{ID: p.ID, Name: p.Name, Description: p.Description}
	}
	return roledto.RoleCredential{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
		Permissions: perms,
	}
}
