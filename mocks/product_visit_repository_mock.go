package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type ProductVisitRepositoryMock struct {
	mock.Mock
}

func NewProductVisitRepositoryMock() *ProductVisitRepositoryMock {
	return &ProductVisitRepositoryMock{}
}

func (m *ProductVisitRepositoryMock) CreateVisit(db database.Database, visit entity.ProductVisit) (*entity.ProductVisit, error) {
	args := m.Called(db, visit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.ProductVisit), args.Error(1)
}

func (m *ProductVisitRepositoryMock) GetVisitCountByProductID(db database.Database, productID uint) (int64, error) {
	args := m.Called(db, productID)
	return int64(args.Int(0)), args.Error(1)
}

func (m *ProductVisitRepositoryMock) HasVisitedInLast24h(db database.Database, productID uint, visitorIP string) (bool, error) {
	args := m.Called(db, productID, visitorIP)
	return args.Bool(0), args.Error(1)
}
