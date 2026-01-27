package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func NewUserRepositoryMock() *UserRepositoryMock {
	return &UserRepositoryMock{}
}

func (u *UserRepositoryMock) FindUsers(db database.Database) ([]*entity.User, error) {
	args := u.Called(db)
	return args.Get(0).([]*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) FindUserByID(db database.Database, id uint) (*entity.User, error) {
	args := u.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) FindUserByPhone(db database.Database, phone string) (*entity.User, error) {
	args := u.Called(db, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) FindUserByEmail(db database.Database, email string) (*entity.User, error) {
	args := u.Called(db, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (u *UserRepositoryMock) CreateUser(db database.Database, user *entity.User) error {
	args := u.Called(db, user)
	return args.Error(0)
}

func (u *UserRepositoryMock) DeleteUserByPhone(db database.Database, phone string) error {
	args := u.Called(db, phone)
	return args.Error(0)
}

func (u *UserRepositoryMock) UpdateUser(db database.Database, user *entity.User) error {
	args := u.Called(db, user)
	return args.Error(0)
}
