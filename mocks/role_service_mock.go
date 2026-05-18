package mocks

import (
	roledto "github.com/Mahoura-shop/Backend/internal/application/dto/role"
	"github.com/stretchr/testify/mock"
)

type RoleServiceMock struct {
	mock.Mock
}

func NewRoleServiceMock() *RoleServiceMock {
	return &RoleServiceMock{}
}

func (m *RoleServiceMock) GetRoles() ([]roledto.RoleCredential, error) {
	args := m.Called()
	return args.Get(0).([]roledto.RoleCredential), args.Error(1)
}

func (m *RoleServiceMock) GetPermissions() ([]roledto.PermissionCredential, error) {
	args := m.Called()
	return args.Get(0).([]roledto.PermissionCredential), args.Error(1)
}

func (m *RoleServiceMock) CreateRole(req roledto.CreateRoleRequest) (*roledto.RoleCredential, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*roledto.RoleCredential), args.Error(1)
}

func (m *RoleServiceMock) UpdateRole(req roledto.UpdateRoleRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *RoleServiceMock) DeleteRole(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
