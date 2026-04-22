package mocks

import (
	userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *UserServiceMock {
	return &UserServiceMock{}
}

func (s *UserServiceMock) IsUserActive(userID uint) error {
	args := s.Called(userID)
	return args.Error(0)
}

func (s *UserServiceMock) GetUserByID(userID uint) (*entity.User, error) {
	args := s.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (s *UserServiceMock) GetUserCredential(userID uint) (userdto.UserCredential, error) {
	args := s.Called(userID)
	return args.Get(0).(userdto.UserCredential), args.Error(1)
}

func (s *UserServiceMock) BanUser(userID uint) error {
	args := s.Called(userID)
	return args.Error(0)
}

func (s *UserServiceMock) UnbanUser(userID uint) error {
	args := s.Called(userID)
	return args.Error(0)
}

func (s *UserServiceMock) VerifyEmail(verifyOTPInfo userdto.VerifyEmailRequest) error {
	args := s.Called(verifyOTPInfo)
	return args.Error(0)
}

func (s *UserServiceMock) FindActiveUserByPhone(phone string) (*entity.User, error) {
	args := s.Called(phone)
	return args.Get(0).(*entity.User), args.Error(1)
}