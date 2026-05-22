package service

import (
	"errors"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/suite"
)

type BrandServiceTestSuite struct {
	suite.Suite
	constants         *bootstrap.Constants
	brandRepository   *mocks.BrandRepositoryMock
	productRepository *mocks.ProductRepositoryMock
	db                *mocks.DatabaseMock
	service           *BrandService
}

func (s *BrandServiceTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.brandRepository = mocks.NewBrandRepositoryMock()
	s.productRepository = mocks.NewProductRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewBrandService(BrandServiceDeps{
		Constants:         s.constants,
		BrandRepository:   s.brandRepository,
		ProductRepository: s.productRepository,
		DB:                s.db,
	})
}

func (s *BrandServiceTestSuite) TestFindBrandByID_Success() {
	brand := &entity.Brand{
		Model:    dbmodel.Model{ID: 5},
		Name:     "Nivea",
		Slug:     "nivea",
		IsActive: true,
	}
	s.brandRepository.On("FindBrandByID", s.db, uint(5)).Return(brand, nil).Once()
	s.brandRepository.On("GetBrandProductsCount", s.db, uint(5)).Return(uint(3), nil).Once()

	result, err := s.service.FindBrandByID(5)

	s.NoError(err)
	s.NotNil(result)
	s.Equal("Nivea", result.Name)
	s.Equal(uint(3), result.Count)
	s.brandRepository.AssertExpectations(s.T())
}

func (s *BrandServiceTestSuite) TestFindBrandByID_NotFound() {
	var nilBrand *entity.Brand
	s.brandRepository.On("FindBrandByID", s.db, uint(99)).Return(nilBrand, nil).Once()

	result, err := s.service.FindBrandByID(99)

	s.Nil(result)
	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.brandRepository.AssertExpectations(s.T())
}

func (s *BrandServiceTestSuite) TestFindBrandByID_RepoError() {
	var nilBrand *entity.Brand
	repoErr := errors.New("db error")
	s.brandRepository.On("FindBrandByID", s.db, uint(1)).Return(nilBrand, repoErr).Once()

	result, err := s.service.FindBrandByID(1)

	s.Nil(result)
	s.ErrorIs(err, repoErr)
	s.brandRepository.AssertExpectations(s.T())
}

func (s *BrandServiceTestSuite) TestFindBrandBySlug_Success() {
	brand := &entity.Brand{
		Model: dbmodel.Model{ID: 7},
		Name:  "Loreal",
		Slug:  "loreal",
	}
	s.brandRepository.On("FindBrandBySlug", s.db, "loreal").Return(brand, nil).Once()
	s.brandRepository.On("GetBrandProductsCount", s.db, uint(7)).Return(uint(0), nil).Once()

	result, err := s.service.FindBrandBySlug("loreal")

	s.NoError(err)
	s.NotNil(result)
	s.Equal("loreal", result.Slug)
	s.brandRepository.AssertExpectations(s.T())
}

func (s *BrandServiceTestSuite) TestFindBrandBySlug_NotFound() {
	var nilBrand *entity.Brand
	s.brandRepository.On("FindBrandBySlug", s.db, "missing").Return(nilBrand, nil).Once()

	result, err := s.service.FindBrandBySlug("missing")

	s.Nil(result)
	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.brandRepository.AssertExpectations(s.T())
}

func (s *BrandServiceTestSuite) TestDeleteBrand_Success() {
	brand := &entity.Brand{Model: dbmodel.Model{ID: 5}}
	s.brandRepository.On("FindBrandByID", s.db, uint(5)).Return(brand, nil).Once()
	s.brandRepository.On("DeleteBrandByID", s.db, uint(5)).Return(nil).Once()

	err := s.service.DeleteBrand(5)

	s.NoError(err)
	s.brandRepository.AssertExpectations(s.T())
}

func (s *BrandServiceTestSuite) TestDeleteBrand_NotFound() {
	var nilBrand *entity.Brand
	s.brandRepository.On("FindBrandByID", s.db, uint(99)).Return(nilBrand, nil).Once()

	err := s.service.DeleteBrand(99)

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.brandRepository.AssertNotCalled(s.T(), "DeleteBrandByID")
	s.brandRepository.AssertExpectations(s.T())
}

func (s *BrandServiceTestSuite) TestParseBrand_FieldMapping() {
	brand := entity.Brand{
		Model:       dbmodel.Model{ID: 3},
		Name:        "Test Brand",
		Slug:        "test-brand",
		Description: "A brand",
		IsActive:    true,
		BrandPic:    "",
	}
	s.brandRepository.On("GetBrandProductsCount", s.db, uint(3)).Return(uint(10), nil).Once()

	result, err := s.service.ParseBrand(brand)

	s.NoError(err)
	s.Equal(uint(3), result.ID)
	s.Equal("Test Brand", result.Name)
	s.Equal("test-brand", result.Slug)
	s.Equal(uint(10), result.Count)
	s.brandRepository.AssertExpectations(s.T())
}

func TestBrandServiceSuite(t *testing.T) {
	suite.Run(t, new(BrandServiceTestSuite))
}
