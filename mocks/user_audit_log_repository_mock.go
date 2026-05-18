package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type UserAuditLogRepositoryMock struct {
	mock.Mock
}

func NewUserAuditLogRepositoryMock() *UserAuditLogRepositoryMock {
	return &UserAuditLogRepositoryMock{}
}

func (m *UserAuditLogRepositoryMock) CreateUserAuditLog(db database.Database, log entity.UserAuditLog) error {
	args := m.Called(db, log)
	return args.Error(0)
}

func (m *UserAuditLogRepositoryMock) GetAuditLogsByUserID(db database.Database, userID uint) ([]*entity.UserAuditLog, error) {
	args := m.Called(db, userID)
	return args.Get(0).([]*entity.UserAuditLog), args.Error(1)
}
