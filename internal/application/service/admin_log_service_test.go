package service

import (
	"errors"
	"testing"
	"time"

	adminlogdto "github.com/Mahoura-shop/Backend/internal/application/dto/admin_log"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/suite"
)

type AdminLogServiceTestSuite struct {
	suite.Suite
	repo    *mocks.AdminActivityLogRepositoryMock
	db      *mocks.DatabaseMock
	service *AdminLogService
}

func (s *AdminLogServiceTestSuite) SetupTest() {
	s.repo = mocks.NewAdminActivityLogRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewAdminLogService(AdminLogServiceDeps{
		Repo: s.repo,
		DB:   s.db,
	})
}

func (s *AdminLogServiceTestSuite) TestCreateLog_Success() {
	log := entity.AdminActivityLog{
		AdminID:    1,
		Method:     "POST",
		Path:       "/admin/products",
		IPAddress:  "127.0.0.1",
		StatusCode: 201,
	}
	s.repo.On("Create", s.db, log).Return(nil).Once()

	err := s.service.CreateLog(1, "POST", "/admin/products", "127.0.0.1", 201)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *AdminLogServiceTestSuite) TestCreateLog_RepoError() {
	log := entity.AdminActivityLog{
		AdminID:    2,
		Method:     "DELETE",
		Path:       "/admin/users/5",
		IPAddress:  "10.0.0.1",
		StatusCode: 200,
	}
	s.repo.On("Create", s.db, log).Return(errors.New("db error")).Once()

	err := s.service.CreateLog(2, "DELETE", "/admin/users/5", "10.0.0.1", 200)

	s.Error(err)
	s.repo.AssertExpectations(s.T())
}

func (s *AdminLogServiceTestSuite) TestGetLogs_ReturnsMappedDTOs() {
	now := time.Now()
	admin := entity.User{
		Model:     dbmodel.Model{ID: 1},
		Phone:     "+989123456789",
		FirstName: "Ali",
		LastName:  "Rezaei",
	}
	logs := []*entity.AdminActivityLog{
		{
			Model:      dbmodel.Model{ID: 10, CreatedAt: now},
			AdminID:    1,
			Admin:      admin,
			Method:     "GET",
			Path:       "/admin/orders",
			IPAddress:  "192.168.1.1",
			StatusCode: 200,
		},
	}
	s.repo.On("GetAll", s.db).Return(logs, nil).Once()

	result, err := s.service.GetLogs()

	s.NoError(err)
	s.Require().Len(result, 1)
	dto := result[0]
	s.Equal(uint(10), dto.ID)
	s.Equal(uint(1), dto.AdminID)
	s.Equal("+989123456789", dto.AdminPhone)
	s.Equal("Ali Rezaei", dto.AdminName)
	s.Equal("GET", dto.Method)
	s.Equal("/admin/orders", dto.Path)
	s.Equal(200, dto.StatusCode)
	s.repo.AssertExpectations(s.T())
}

func (s *AdminLogServiceTestSuite) TestGetLogs_RepoError() {
	s.repo.On("GetAll", s.db).Return(nil, errors.New("db error")).Once()

	result, err := s.service.GetLogs()

	s.Nil(result)
	s.Error(err)
	s.repo.AssertExpectations(s.T())
}

func (s *AdminLogServiceTestSuite) TestGetLogs_Empty() {
	s.repo.On("GetAll", s.db).Return([]*entity.AdminActivityLog{}, nil).Once()

	result, err := s.service.GetLogs()

	s.NoError(err)
	s.Empty(result)
	s.repo.AssertExpectations(s.T())
}

func (s *AdminLogServiceTestSuite) TestGetLogs_ReturnsCorrectType() {
	s.repo.On("GetAll", s.db).Return([]*entity.AdminActivityLog{}, nil).Once()

	result, err := s.service.GetLogs()

	s.NoError(err)
	_, ok := interface{}(result).([]adminlogdto.AdminActivityLogDTO)
	s.True(ok)
	s.repo.AssertExpectations(s.T())
}

func TestAdminLogServiceSuite(t *testing.T) {
	suite.Run(t, new(AdminLogServiceTestSuite))
}
