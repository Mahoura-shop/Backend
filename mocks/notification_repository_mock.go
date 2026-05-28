package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type NotificationRepositoryMock struct {
	mock.Mock
}

func NewNotificationRepositoryMock() *NotificationRepositoryMock {
	return &NotificationRepositoryMock{}
}

func (m *NotificationRepositoryMock) Create(db database.Database, n entity.Notification) (*entity.Notification, error) {
	args := m.Called(db, n)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Notification), args.Error(1)
}

func (m *NotificationRepositoryMock) GetByUserID(db database.Database, userID uint) ([]*entity.Notification, error) {
	args := m.Called(db, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Notification), args.Error(1)
}

func (m *NotificationRepositoryMock) MarkAsRead(db database.Database, id, userID uint) error {
	args := m.Called(db, id, userID)
	return args.Error(0)
}

func (m *NotificationRepositoryMock) MarkAllAsRead(db database.Database, userID uint) error {
	args := m.Called(db, userID)
	return args.Error(0)
}

func (m *NotificationRepositoryMock) CountUnread(db database.Database, userID uint) (int64, error) {
	args := m.Called(db, userID)
	return int64(args.Int(0)), args.Error(1)
}
