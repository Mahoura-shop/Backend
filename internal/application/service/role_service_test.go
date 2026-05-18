package service

import (
	"errors"
	"testing"

	roledto "github.com/Mahoura-shop/Backend/internal/application/dto/role"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RoleServiceTestSuite struct {
	suite.Suite
	roleRepository *mocks.RoleRepositoryMock
	db             *mocks.DatabaseMock
	service        *RoleService
}

func (s *RoleServiceTestSuite) SetupTest() {
	s.roleRepository = mocks.NewRoleRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewRoleService(RoleServiceDeps{
		RoleRepository: s.roleRepository,
		DB:             s.db,
	})
}

func (s *RoleServiceTestSuite) TestGetRoles_Success() {
	roles := []*entity.Role{
		{Model: dbmodel.Model{ID: 1}, Name: "Admin", Permissions: []entity.Permission{}},
		{Model: dbmodel.Model{ID: 2}, Name: "Editor", Permissions: []entity.Permission{}},
	}
	s.roleRepository.On("GetRoles", s.db).Return(roles, nil).Once()

	result, err := s.service.GetRoles()

	s.NoError(err)
	s.Len(result, 2)
	s.Equal("Admin", result[0].Name)
	s.roleRepository.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestGetRoles_Empty() {
	s.roleRepository.On("GetRoles", s.db).Return([]*entity.Role{}, nil).Once()

	result, err := s.service.GetRoles()

	s.NoError(err)
	s.Empty(result)
	s.roleRepository.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestGetPermissions_Success() {
	perms := []*entity.Permission{
		{Model: dbmodel.Model{ID: 1}, Name: "read", Description: "Read access"},
		{Model: dbmodel.Model{ID: 2}, Name: "write", Description: "Write access"},
	}
	s.roleRepository.On("GetPermissions", s.db).Return(perms, nil).Once()

	result, err := s.service.GetPermissions()

	s.NoError(err)
	s.Len(result, 2)
	s.Equal("read", result[0].Name)
	s.roleRepository.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestCreateRole_Success() {
	req := roledto.CreateRoleRequest{
		Name:          "Moderator",
		Description:   "Moderate content",
		PermissionIDs: []uint{1, 2},
	}
	createdRole := &entity.Role{
		Model:       dbmodel.Model{ID: 3},
		Name:        "Moderator",
		Description: "Moderate content",
		Permissions: []entity.Permission{
			{Model: dbmodel.Model{ID: 1}},
			{Model: dbmodel.Model{ID: 2}},
		},
	}
	s.roleRepository.On("CreateRole", s.db, mock.MatchedBy(func(r entity.Role) bool {
		return r.Name == "Moderator" && r.Description == "Moderate content" && len(r.Permissions) == 2
	})).Return(createdRole, nil).Once()

	result, err := s.service.CreateRole(req)

	s.NoError(err)
	s.NotNil(result)
	s.Equal("Moderator", result.Name)
	s.roleRepository.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestCreateRole_RepoError() {
	req := roledto.CreateRoleRequest{Name: "Bad", PermissionIDs: []uint{}}
	repoErr := errors.New("db error")
	s.roleRepository.On("CreateRole", s.db, mock.MatchedBy(func(r entity.Role) bool {
		return r.Name == "Bad"
	})).Return((*entity.Role)(nil), repoErr).Once()

	result, err := s.service.CreateRole(req)

	s.Nil(result)
	s.ErrorIs(err, repoErr)
	s.roleRepository.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestDeleteRole_CallsRepo() {
	s.roleRepository.On("DeleteRoleByID", s.db, uint(5)).Return(nil).Once()

	err := s.service.DeleteRole(5)

	s.NoError(err)
	s.roleRepository.AssertExpectations(s.T())
}

func (s *RoleServiceTestSuite) TestMapRole_PermissionsAreMapped() {
	role := &entity.Role{
		Model: dbmodel.Model{ID: 1},
		Name:  "Admin",
		Permissions: []entity.Permission{
			{Model: dbmodel.Model{ID: 10}, Name: "create", Description: "Create"},
		},
	}
	s.roleRepository.On("GetRoles", s.db).Return([]*entity.Role{role}, nil).Once()

	result, err := s.service.GetRoles()

	s.NoError(err)
	s.Len(result[0].Permissions, 1)
	s.Equal(uint(10), result[0].Permissions[0].ID)
	s.roleRepository.AssertExpectations(s.T())
}

func TestRoleServiceSuite(t *testing.T) {
	suite.Run(t, new(RoleServiceTestSuite))
}
