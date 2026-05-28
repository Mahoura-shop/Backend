package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	postgresrepo "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
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

func (m *ProductVisitRepositoryMock) GetVisitsPerDay(db database.Database, productID uint, days int) ([]postgresrepo.VisitsByDay, error) {
	args := m.Called(db, productID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgresrepo.VisitsByDay), args.Error(1)
}

func (m *ProductVisitRepositoryMock) GetAllVisitsPerDay(db database.Database, days int) ([]postgresrepo.VisitsByDay, error) {
	args := m.Called(db, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgresrepo.VisitsByDay), args.Error(1)
}

func (m *ProductVisitRepositoryMock) GetCategoryVisitsPerDay(db database.Database, categoryID uint, days int) ([]postgresrepo.VisitsByDay, error) {
	args := m.Called(db, categoryID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgresrepo.VisitsByDay), args.Error(1)
}

func (m *ProductVisitRepositoryMock) GetBrandVisitsPerDay(db database.Database, brandID uint, days int) ([]postgresrepo.VisitsByDay, error) {
	args := m.Called(db, brandID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgresrepo.VisitsByDay), args.Error(1)
}
