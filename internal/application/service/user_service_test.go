package service

import (
	"errors"
	"testing"

	userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"
	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/suite"
)

type UserServiceTestSuite struct {
	suite.Suite
	constants  *bootstrap.Constants
	userRepo   *mocks.UserRepositoryMock
	db         *mocks.DatabaseMock
	service    *UserService
}

func (s *UserServiceTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.userRepo = mocks.NewUserRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewUserService(UserServiceDeps{
		Constants:      s.constants,
		UserRepository: s.userRepo,
		DB:             s.db,
	})
}

func sampleUser() entity.User {
	return entity.User{
		Model:     dbmodel.Model{ID: 1},
		FirstName: "Ali",
		LastName:  "Rezaei",
		Phone:     "+989123456789",
		Email:     "ali@example.com",
		Status:    enum.UserStatusActive,
	}
}

func (s *UserServiceTestSuite) TestParseUser_MapsBasicFields() {
	user := sampleUser()

	result := s.service.ParseUser(user)

	s.Equal(uint(1), result.ID)
	s.Equal("Ali", result.FirstName)
	s.Equal("Rezaei", result.LastName)
	s.Equal("+989123456789", result.Phone)
	s.Equal("ali@example.com", result.Email)
}

func (s *UserServiceTestSuite) TestParseUser_WithRole_PopulatesRoleName() {
	user := sampleUser()
	user.Role = &entity.Role{Name: "admin"}

	result := s.service.ParseUser(user)

	s.Equal("admin", result.RoleName)
}

func (s *UserServiceTestSuite) TestParseUser_NoRole_EmptyRoleName() {
	user := sampleUser()

	result := s.service.ParseUser(user)

	s.Equal("", result.RoleName)
}

func (s *UserServiceTestSuite) TestParseUser_ReturnsCorrectType() {
	user := sampleUser()

	result := s.service.ParseUser(user)

	_, ok := interface{}(result).(userdto.UserCredential)
	s.True(ok)
}

func (s *UserServiceTestSuite) TestGetUserByID_Success() {
	user := sampleUser()
	s.userRepo.On("FindUserByID", s.db, uint(1)).Return(&user, nil).Once()

	result, err := s.service.GetUserByID(1)

	s.NoError(err)
	s.Require().NotNil(result)
	s.Equal(uint(1), result.ID)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestGetUserByID_NotFound() {
	var nilUser *entity.User
	s.userRepo.On("FindUserByID", s.db, uint(99)).Return(nilUser, nil).Once()

	result, err := s.service.GetUserByID(99)

	s.Nil(result)
	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
}

func (s *UserServiceTestSuite) TestGetUserByID_RepoError() {
	s.userRepo.On("FindUserByID", s.db, uint(1)).Return(nil, errors.New("db error")).Once()

	result, err := s.service.GetUserByID(1)

	s.Nil(result)
	s.Error(err)
}

func (s *UserServiceTestSuite) TestGetUsers_ReturnsList() {
	user := sampleUser()
	s.userRepo.On("FindUsers", s.db).Return([]*entity.User{&user}, nil).Once()

	result, err := s.service.GetUsers()

	s.NoError(err)
	s.Len(result, 1)
	s.Equal("+989123456789", result[0].Phone)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestGetUsers_Empty() {
	s.userRepo.On("FindUsers", s.db).Return([]*entity.User{}, nil).Once()

	result, err := s.service.GetUsers()

	s.NoError(err)
	s.Empty(result)
}

func (s *UserServiceTestSuite) TestGetUsers_RepoError() {
	s.userRepo.On("FindUsers", s.db).Return(nil, errors.New("db error")).Once()

	_, err := s.service.GetUsers()

	s.Error(err)
}

func (s *UserServiceTestSuite) TestBanUser_Success() {
	user := sampleUser()
	s.userRepo.On("FindUserByID", s.db, uint(1)).Return(&user, nil).Once()
	bannedUser := sampleUser()
	bannedUser.Status = enum.UserStatusBlock
	s.userRepo.On("UpdateUser", s.db, bannedUser).Return(nil).Once()

	err := s.service.BanUser(1)

	s.NoError(err)
	s.userRepo.AssertExpectations(s.T())
}

func (s *UserServiceTestSuite) TestBanUser_AlreadyBanned_ConflictError() {
	user := sampleUser()
	user.Status = enum.UserStatusBlock
	s.userRepo.On("FindUserByID", s.db, uint(1)).Return(&user, nil).Once()

	err := s.service.BanUser(1)

	s.Error(err)
	var ce exception.ConflictErrors
	s.True(errors.As(err, &ce))
	s.userRepo.AssertNotCalled(s.T(), "UpdateUser")
}

func (s *UserServiceTestSuite) TestBanUser_NotFound() {
	var nilUser *entity.User
	s.userRepo.On("FindUserByID", s.db, uint(99)).Return(nilUser, nil).Once()

	err := s.service.BanUser(99)

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
}

func (s *UserServiceTestSuite) TestIsUserActive_ActiveUser_NoError() {
	user := sampleUser()
	s.userRepo.On("FindUserByID", s.db, uint(1)).Return(&user, nil).Once()

	err := s.service.IsUserActive(1)

	s.NoError(err)
}

func (s *UserServiceTestSuite) TestIsUserActive_BlockedUser_ReturnsError() {
	user := sampleUser()
	user.Status = enum.UserStatusBlock
	s.userRepo.On("FindUserByID", s.db, uint(1)).Return(&user, nil).Once()

	err := s.service.IsUserActive(1)

	s.Error(err)
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}
