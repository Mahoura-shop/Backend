package mocks

import (
	categorydto "github.com/Mahoura-shop/Backend/internal/application/dto/category"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type CategoryServiceMock struct {
	mock.Mock
}

func NewCategoryServiceMock() *CategoryServiceMock {
	return &CategoryServiceMock{}
}

func (m *CategoryServiceMock) ParseCategory(c entity.Category) (categorydto.CategoryCredential, error) {
	args := m.Called(c)
	return args.Get(0).(categorydto.CategoryCredential), args.Error(1)
}

func (m *CategoryServiceMock) FindCategoryBySlug(slug string) (*categorydto.CategoryCredential, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*categorydto.CategoryCredential), args.Error(1)
}

func (m *CategoryServiceMock) FindCategoryByID(id uint) (*categorydto.CategoryCredential, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*categorydto.CategoryCredential), args.Error(1)
}

func (m *CategoryServiceMock) CreateCategory(req categorydto.CreateCategoryRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *CategoryServiceMock) GetCategories() ([]categorydto.CategoryCredential, error) {
	args := m.Called()
	return args.Get(0).([]categorydto.CategoryCredential), args.Error(1)
}

func (m *CategoryServiceMock) DeleteCategory(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *CategoryServiceMock) UpdateCategory(req categorydto.UpdateCategoryRequest) error {
	args := m.Called(req)
	return args.Error(0)
}
