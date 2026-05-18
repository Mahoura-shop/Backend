package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type RoleRepositoryMock struct {
	mock.Mock
}

func NewRoleRepositoryMock() *RoleRepositoryMock {
	return &RoleRepositoryMock{}
}

func (m *RoleRepositoryMock) GetRoles(db database.Database) ([]*entity.Role, error) {
	args := m.Called(db)
	return args.Get(0).([]*entity.Role), args.Error(1)
}

func (m *RoleRepositoryMock) FindRoleByID(db database.Database, id uint) (*entity.Role, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *RoleRepositoryMock) CreateRole(db database.Database, role entity.Role) (*entity.Role, error) {
	args := m.Called(db, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Role), args.Error(1)
}

func (m *RoleRepositoryMock) UpdateRole(db database.Database, role entity.Role) error {
	args := m.Called(db, role)
	return args.Error(0)
}

func (m *RoleRepositoryMock) DeleteRoleByID(db database.Database, id uint) error {
	args := m.Called(db, id)
	return args.Error(0)
}

func (m *RoleRepositoryMock) GetPermissions(db database.Database) ([]*entity.Permission, error) {
	args := m.Called(db)
	return args.Get(0).([]*entity.Permission), args.Error(1)
}
