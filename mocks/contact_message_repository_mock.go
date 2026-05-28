package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type ContactMessageRepositoryMock struct {
	mock.Mock
}

func NewContactMessageRepositoryMock() *ContactMessageRepositoryMock {
	return &ContactMessageRepositoryMock{}
}

func (m *ContactMessageRepositoryMock) Create(db database.Database, msg entity.ContactMessage) (*entity.ContactMessage, error) {
	args := m.Called(db, msg)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ContactMessage), args.Error(1)
}

func (m *ContactMessageRepositoryMock) GetAll(db database.Database) ([]*entity.ContactMessage, error) {
	args := m.Called(db)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.ContactMessage), args.Error(1)
}
