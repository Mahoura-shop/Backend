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

type CategoryServiceTestSuite struct {
	suite.Suite
	constants          *bootstrap.Constants
	categoryRepository *mocks.CategoryRepositoryMock
	db                 *mocks.DatabaseMock
	service            *CategoryService
}

func (s *CategoryServiceTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.categoryRepository = mocks.NewCategoryRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewCategoryService(CategoryServiceDeps{
		Constants:          s.constants,
		CategoryRepository: s.categoryRepository,
		DB:                 s.db,
	})
}

func (s *CategoryServiceTestSuite) TestFindCategoryByID_Success() {
	category := &entity.Category{
		Model:    dbmodel.Model{ID: 3},
		Name:     "Skincare",
		Slug:     "skincare",
		IsActive: true,
	}
	s.categoryRepository.On("FindCategoryByID", s.db, uint(3)).Return(category, nil).Once()
	s.categoryRepository.On("GetCategoryProductsCount", s.db, uint(3)).Return(uint(12), nil).Once()

	result, err := s.service.FindCategoryByID(3)

	s.NoError(err)
	s.NotNil(result)
	s.Equal("Skincare", result.Name)
	s.Equal(uint(12), result.Count)
	s.categoryRepository.AssertExpectations(s.T())
}

func (s *CategoryServiceTestSuite) TestFindCategoryByID_NotFound() {
	var nilCategory *entity.Category
	s.categoryRepository.On("FindCategoryByID", s.db, uint(99)).Return(nilCategory, nil).Once()

	result, err := s.service.FindCategoryByID(99)

	s.Nil(result)
	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.categoryRepository.AssertExpectations(s.T())
}

func (s *CategoryServiceTestSuite) TestFindCategoryByID_RepoError() {
	var nilCategory *entity.Category
	repoErr := errors.New("db error")
	s.categoryRepository.On("FindCategoryByID", s.db, uint(1)).Return(nilCategory, repoErr).Once()

	result, err := s.service.FindCategoryByID(1)

	s.Nil(result)
	s.ErrorIs(err, repoErr)
	s.categoryRepository.AssertExpectations(s.T())
}

func (s *CategoryServiceTestSuite) TestFindCategoryBySlug_Success() {
	category := &entity.Category{
		Model: dbmodel.Model{ID: 4},
		Name:  "Makeup",
		Slug:  "makeup",
	}
	s.categoryRepository.On("FindCategoryBySlug", s.db, "makeup").Return(category, nil).Once()
	s.categoryRepository.On("GetCategoryProductsCount", s.db, uint(4)).Return(uint(5), nil).Once()

	result, err := s.service.FindCategoryBySlug("makeup")

	s.NoError(err)
	s.NotNil(result)
	s.Equal("makeup", result.Slug)
	s.categoryRepository.AssertExpectations(s.T())
}

func (s *CategoryServiceTestSuite) TestFindCategoryBySlug_NotFound() {
	var nilCategory *entity.Category
	s.categoryRepository.On("FindCategoryBySlug", s.db, "missing").Return(nilCategory, nil).Once()

	result, err := s.service.FindCategoryBySlug("missing")

	s.Nil(result)
	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.categoryRepository.AssertExpectations(s.T())
}

func (s *CategoryServiceTestSuite) TestDeleteCategory_Success() {
	category := &entity.Category{Model: dbmodel.Model{ID: 3}}
	s.categoryRepository.On("FindCategoryByID", s.db, uint(3)).Return(category, nil).Once()
	s.categoryRepository.On("DeleteCategoryByID", s.db, uint(3)).Return(nil).Once()

	err := s.service.DeleteCategory(3)

	s.NoError(err)
	s.categoryRepository.AssertExpectations(s.T())
}

func (s *CategoryServiceTestSuite) TestDeleteCategory_NotFound() {
	var nilCategory *entity.Category
	s.categoryRepository.On("FindCategoryByID", s.db, uint(99)).Return(nilCategory, nil).Once()

	err := s.service.DeleteCategory(99)

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.categoryRepository.AssertNotCalled(s.T(), "DeleteCategoryByID")
	s.categoryRepository.AssertExpectations(s.T())
}

func (s *CategoryServiceTestSuite) TestGetCategoryProductsCount_ReturnsZeroOnError() {
	repoErr := errors.New("count error")
	s.categoryRepository.On("GetCategoryProductsCount", s.db, uint(1)).Return(uint(0), repoErr).Once()

	count, err := s.service.GetCategoryProductsCount(1)

	s.Error(err)
	s.Equal(uint(0), count)
	s.categoryRepository.AssertExpectations(s.T())
}

func TestCategoryServiceSuite(t *testing.T) {
	suite.Run(t, new(CategoryServiceTestSuite))
}
