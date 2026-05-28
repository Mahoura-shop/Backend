package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type ReturnRepositoryMock struct {
	mock.Mock
}

func NewReturnRepositoryMock() *ReturnRepositoryMock {
	return &ReturnRepositoryMock{}
}

func (m *ReturnRepositoryMock) CreateReturn(db database.Database, r entity.Return) (*entity.Return, error) {
	args := m.Called(db, r)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Return), args.Error(1)
}

func (m *ReturnRepositoryMock) FindReturnByID(db database.Database, id uint) (*entity.Return, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Return), args.Error(1)
}

func (m *ReturnRepositoryMock) GetReturnsByUserID(db database.Database, userID uint) ([]*entity.Return, error) {
	args := m.Called(db, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Return), args.Error(1)
}

func (m *ReturnRepositoryMock) GetAllReturns(db database.Database) ([]*entity.Return, error) {
	args := m.Called(db)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Return), args.Error(1)
}

func (m *ReturnRepositoryMock) GetReturnsByStatus(db database.Database, status enum.ReturnStatus) ([]*entity.Return, error) {
	args := m.Called(db, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Return), args.Error(1)
}

func (m *ReturnRepositoryMock) UpdateReturn(db database.Database, r entity.Return) error {
	args := m.Called(db, r)
	return args.Error(0)
}
