package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type UpgradeRequestRepositoryMock struct {
	mock.Mock
}

func NewUpgradeRequestRepositoryMock() *UpgradeRequestRepositoryMock {
	return &UpgradeRequestRepositoryMock{}
}

func (m *UpgradeRequestRepositoryMock) CreateUpgradeRequest(db database.Database, req entity.UpgradeRequest) (*entity.UpgradeRequest, error) {
	args := m.Called(db, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UpgradeRequest), args.Error(1)
}

func (m *UpgradeRequestRepositoryMock) FindUpgradeRequestByID(db database.Database, id uint) (*entity.UpgradeRequest, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UpgradeRequest), args.Error(1)
}

func (m *UpgradeRequestRepositoryMock) GetUpgradeRequests(db database.Database) ([]*entity.UpgradeRequest, error) {
	args := m.Called(db)
	return args.Get(0).([]*entity.UpgradeRequest), args.Error(1)
}

func (m *UpgradeRequestRepositoryMock) GetUpgradeRequestsByUserID(db database.Database, userID uint) ([]*entity.UpgradeRequest, error) {
	args := m.Called(db, userID)
	return args.Get(0).([]*entity.UpgradeRequest), args.Error(1)
}

func (m *UpgradeRequestRepositoryMock) GetUpgradeRequestsByStatus(db database.Database, status enum.UpgradeRequestStatus) ([]*entity.UpgradeRequest, error) {
	args := m.Called(db, status)
	return args.Get(0).([]*entity.UpgradeRequest), args.Error(1)
}

func (m *UpgradeRequestRepositoryMock) FindPendingByUserID(db database.Database, userID uint) (*entity.UpgradeRequest, error) {
	args := m.Called(db, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UpgradeRequest), args.Error(1)
}

func (m *UpgradeRequestRepositoryMock) UpdateUpgradeRequest(db database.Database, req entity.UpgradeRequest) error {
	args := m.Called(db, req)
	return args.Error(0)
}
