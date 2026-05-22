package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	domainPostgres "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type ProductRepositoryMock struct {
	mock.Mock
}

func NewProductRepositoryMock() *ProductRepositoryMock {
	return &ProductRepositoryMock{}
}

func (m *ProductRepositoryMock) CreateProduct(db database.Database, product entity.Product) (*entity.Product, error) {
	args := m.Called(db, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) FindProductByID(db database.Database, id uint) (*entity.Product, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) FindProductBySlug(db database.Database, slug string) (*entity.Product, error) {
	args := m.Called(db, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) FindProductByName(db database.Database, name string) (*entity.Product, error) {
	args := m.Called(db, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetProducts(db database.Database) ([]*entity.Product, error) {
	args := m.Called(db)
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) SearchProducts(db database.Database, filter domainPostgres.ProductFilter) ([]*entity.Product, error) {
	args := m.Called(db, filter)
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) SearchProductsWithCount(db database.Database, filter domainPostgres.ProductFilter) ([]*entity.Product, int64, error) {
	args := m.Called(db, filter)
	return args.Get(0).([]*entity.Product), args.Get(1).(int64), args.Error(2)
}

func (m *ProductRepositoryMock) GetRelatedProducts(db database.Database, productID uint, categoryID uint, brandID uint, limit int) ([]*entity.Product, error) {
	args := m.Called(db, productID, categoryID, brandID, limit)
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetCategoryProducts(db database.Database, categoryID uint) ([]*entity.Product, error) {
	args := m.Called(db, categoryID)
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetBrandProducts(db database.Database, brandID uint) ([]*entity.Product, error) {
	args := m.Called(db, brandID)
	return args.Get(0).([]*entity.Product), args.Error(1)
}

func (m *ProductRepositoryMock) GetProductsCount(db database.Database) (uint, error) {
	args := m.Called(db)
	return args.Get(0).(uint), args.Error(1)
}

func (m *ProductRepositoryMock) DeleteProductByID(db database.Database, id uint) error {
	args := m.Called(db, id)
	return args.Error(0)
}

func (m *ProductRepositoryMock) UpdateProduct(db database.Database, product entity.Product) error {
	args := m.Called(db, product)
	return args.Error(0)
}

func (m *ProductRepositoryMock) GetLowStockProducts(db database.Database, limit int) ([]domainPostgres.LowStockProduct, error) {
	args := m.Called(db, limit)
	return args.Get(0).([]domainPostgres.LowStockProduct), args.Error(1)
}

func (m *ProductRepositoryMock) GetTopSoldProducts(db database.Database, limit int) ([]domainPostgres.TopSoldProduct, error) {
	args := m.Called(db, limit)
	return args.Get(0).([]domainPostgres.TopSoldProduct), args.Error(1)
}
