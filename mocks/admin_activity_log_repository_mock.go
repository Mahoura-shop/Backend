package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type AdminActivityLogRepositoryMock struct {
	mock.Mock
}

func NewAdminActivityLogRepositoryMock() *AdminActivityLogRepositoryMock {
	return &AdminActivityLogRepositoryMock{}
}

func (m *AdminActivityLogRepositoryMock) Create(db database.Database, log entity.AdminActivityLog) error {
	args := m.Called(db, log)
	return args.Error(0)
}

func (m *AdminActivityLogRepositoryMock) GetAll(db database.Database) ([]*entity.AdminActivityLog, error) {
	args := m.Called(db)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.AdminActivityLog), args.Error(1)
}
