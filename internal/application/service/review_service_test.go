package service

import (
	"errors"
	"testing"
	"time"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/suite"
)

type ReviewServiceTestSuite struct {
	suite.Suite
	reviewRepository *mocks.ReviewRepositoryMock
	db               *mocks.DatabaseMock
	service          *ReviewService
}

func (s *ReviewServiceTestSuite) SetupTest() {
	s.reviewRepository = mocks.NewReviewRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewReviewService(ReviewServiceDeps{
		ReviewRepository: s.reviewRepository,
		DB:               s.db,
	})
}

func (s *ReviewServiceTestSuite) TestGetProductReviews_ReturnsMapped() {
	now := time.Now()
	reviews := []*entity.Review{
		{
			Model:      dbmodel.Model{ID: 1},
			UserID:     5,
			ProductID:  10,
			Rating:     4,
			Comment:    "Great product",
			IsVerified: true,
			User:       entity.User{Model: dbmodel.Model{ID: 5}, Phone: "09120000000"},
		},
	}
	reviews[0].CreatedAt = now

	s.reviewRepository.On("GetProductReviews", s.db, uint(10)).Return(reviews, nil).Once()

	result, err := s.service.GetProductReviews(10)

	s.NoError(err)
	s.Len(result, 1)
	s.Equal(uint(4), result[0].Rating)
	s.Equal("Great product", result[0].Comment)
	s.True(result[0].IsVerified)
	s.Equal("09120000000", result[0].UserPhone)
	s.reviewRepository.AssertExpectations(s.T())
}

func (s *ReviewServiceTestSuite) TestGetProductReviews_Empty() {
	s.reviewRepository.On("GetProductReviews", s.db, uint(99)).Return([]*entity.Review{}, nil).Once()

	result, err := s.service.GetProductReviews(99)

	s.NoError(err)
	s.Nil(result)
	s.reviewRepository.AssertExpectations(s.T())
}

func (s *ReviewServiceTestSuite) TestGetProductReviews_RepoError() {
	repoErr := errors.New("db error")
	s.reviewRepository.On("GetProductReviews", s.db, uint(10)).Return([]*entity.Review{}, repoErr).Once()

	result, err := s.service.GetProductReviews(10)

	s.Nil(result)
	s.ErrorIs(err, repoErr)
	s.reviewRepository.AssertExpectations(s.T())
}

func (s *ReviewServiceTestSuite) TestDeleteReview_Success() {
	review := &entity.Review{Model: dbmodel.Model{ID: 3}}
	s.reviewRepository.On("FindReviewByID", s.db, uint(3)).Return(review, nil).Once()
	s.reviewRepository.On("DeleteReviewByID", s.db, uint(3)).Return(nil).Once()

	err := s.service.DeleteReview(3)

	s.NoError(err)
	s.reviewRepository.AssertExpectations(s.T())
}

func (s *ReviewServiceTestSuite) TestDeleteReview_NotFound() {
	var nilReview *entity.Review
	s.reviewRepository.On("FindReviewByID", s.db, uint(99)).Return(nilReview, nil).Once()

	err := s.service.DeleteReview(99)

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.reviewRepository.AssertNotCalled(s.T(), "DeleteReviewByID")
	s.reviewRepository.AssertExpectations(s.T())
}

func (s *ReviewServiceTestSuite) TestGetProductReviews_UserPhoneEmptyWhenUserIDZero() {
	reviews := []*entity.Review{
		{
			Model:     dbmodel.Model{ID: 2},
			ProductID: 5,
			Rating:    3,
			User:      entity.User{},
		},
	}
	s.reviewRepository.On("GetProductReviews", s.db, uint(5)).Return(reviews, nil).Once()

	result, err := s.service.GetProductReviews(5)

	s.NoError(err)
	s.Equal("", result[0].UserPhone)
	s.reviewRepository.AssertExpectations(s.T())
}

func TestReviewServiceSuite(t *testing.T) {
	suite.Run(t, new(ReviewServiceTestSuite))
}
