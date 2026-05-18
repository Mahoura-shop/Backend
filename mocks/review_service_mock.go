package mocks

import (
	reviewdto "github.com/Mahoura-shop/Backend/internal/application/dto/review"
	"github.com/stretchr/testify/mock"
)

type ReviewServiceMock struct {
	mock.Mock
}

func NewReviewServiceMock() *ReviewServiceMock {
	return &ReviewServiceMock{}
}

func (m *ReviewServiceMock) SubmitReview(req reviewdto.SubmitReviewRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *ReviewServiceMock) GetProductReviews(productID uint) ([]reviewdto.ReviewCredential, error) {
	args := m.Called(productID)
	return args.Get(0).([]reviewdto.ReviewCredential), args.Error(1)
}

func (m *ReviewServiceMock) GetMyReviews(userID uint) ([]reviewdto.ReviewCredential, error) {
	args := m.Called(userID)
	return args.Get(0).([]reviewdto.ReviewCredential), args.Error(1)
}

func (m *ReviewServiceMock) DeleteReview(reviewID uint) error {
	args := m.Called(reviewID)
	return args.Error(0)
}
