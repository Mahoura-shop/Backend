package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type BrandRepositoryMock struct {
	mock.Mock
}

func NewBrandRepositoryMock() *BrandRepositoryMock {
	return &BrandRepositoryMock{}
}

func (m *BrandRepositoryMock) GetBrandProductsCount(db database.Database, brandID uint) (uint, error) {
	args := m.Called(db, brandID)
	return args.Get(0).(uint), args.Error(1)
}

func (m *BrandRepositoryMock) CreateBrand(db database.Database, brand entity.Brand) (*entity.Brand, error) {
	args := m.Called(db, brand)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Brand), args.Error(1)
}

func (m *BrandRepositoryMock) FindBrandByID(db database.Database, id uint) (*entity.Brand, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Brand), args.Error(1)
}

func (m *BrandRepositoryMock) FindBrandBySlug(db database.Database, slug string) (*entity.Brand, error) {
	args := m.Called(db, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Brand), args.Error(1)
}

func (m *BrandRepositoryMock) GetBrands(db database.Database) ([]*entity.Brand, error) {
	args := m.Called(db)
	return args.Get(0).([]*entity.Brand), args.Error(1)
}

func (m *BrandRepositoryMock) GetBrandsCount(db database.Database) (uint, error) {
	args := m.Called(db)
	return args.Get(0).(uint), args.Error(1)
}

func (m *BrandRepositoryMock) DeleteBrandByID(db database.Database, id uint) error {
	args := m.Called(db, id)
	return args.Error(0)
}

func (m *BrandRepositoryMock) UpdateBrand(db database.Database, brand entity.Brand) error {
	args := m.Called(db, brand)
	return args.Error(0)
}
