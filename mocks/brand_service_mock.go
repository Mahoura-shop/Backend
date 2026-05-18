package mocks

import (
	branddto "github.com/Mahoura-shop/Backend/internal/application/dto/brand"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type BrandServiceMock struct {
	mock.Mock
}

func NewBrandServiceMock() *BrandServiceMock {
	return &BrandServiceMock{}
}

func (m *BrandServiceMock) ParseBrand(b entity.Brand) (branddto.BrandCredential, error) {
	args := m.Called(b)
	return args.Get(0).(branddto.BrandCredential), args.Error(1)
}

func (m *BrandServiceMock) FindBrandBySlug(slug string) (*branddto.BrandCredential, error) {
	args := m.Called(slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*branddto.BrandCredential), args.Error(1)
}

func (m *BrandServiceMock) FindBrandByID(id uint) (*branddto.BrandCredential, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*branddto.BrandCredential), args.Error(1)
}

func (m *BrandServiceMock) CreateBrand(req branddto.CreateBrandRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *BrandServiceMock) GetBrands() ([]branddto.BrandCredential, error) {
	args := m.Called()
	return args.Get(0).([]branddto.BrandCredential), args.Error(1)
}

func (m *BrandServiceMock) DeleteBrand(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *BrandServiceMock) UpdateBrand(req branddto.UpdateBrandRequest) error {
	args := m.Called(req)
	return args.Error(0)
}
