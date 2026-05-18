package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type CategoryRepositoryMock struct {
	mock.Mock
}

func NewCategoryRepositoryMock() *CategoryRepositoryMock {
	return &CategoryRepositoryMock{}
}

func (m *CategoryRepositoryMock) FindCategoryByID(db database.Database, id uint) (*entity.Category, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) FindCategoryBySlug(db database.Database, slug string) (*entity.Category, error) {
	args := m.Called(db, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) CreateCategory(db database.Database, category entity.Category) (*entity.Category, error) {
	args := m.Called(db, category)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) GetCategories(db database.Database) ([]*entity.Category, error) {
	args := m.Called(db)
	return args.Get(0).([]*entity.Category), args.Error(1)
}

func (m *CategoryRepositoryMock) GetCategoriesCount(db database.Database) (uint, error) {
	args := m.Called(db)
	return args.Get(0).(uint), args.Error(1)
}

func (m *CategoryRepositoryMock) DeleteCategoryByID(db database.Database, id uint) error {
	args := m.Called(db, id)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) UpdateCategory(db database.Database, category entity.Category) error {
	args := m.Called(db, category)
	return args.Error(0)
}

func (m *CategoryRepositoryMock) GetCategoryProductsCount(db database.Database, categoryID uint) (uint, error) {
	args := m.Called(db, categoryID)
	return args.Get(0).(uint), args.Error(1)
}
