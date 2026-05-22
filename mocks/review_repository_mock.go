package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type ReviewRepositoryMock struct {
	mock.Mock
}

func NewReviewRepositoryMock() *ReviewRepositoryMock {
	return &ReviewRepositoryMock{}
}

func (m *ReviewRepositoryMock) CreateReview(db database.Database, review entity.Review) (*entity.Review, error) {
	args := m.Called(db, review)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Review), args.Error(1)
}

func (m *ReviewRepositoryMock) FindReviewByID(db database.Database, id uint) (*entity.Review, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Review), args.Error(1)
}

func (m *ReviewRepositoryMock) FindReviewByUserAndProduct(db database.Database, userID, productID uint) (*entity.Review, error) {
	args := m.Called(db, userID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Review), args.Error(1)
}

func (m *ReviewRepositoryMock) GetProductReviews(db database.Database, productID uint) ([]*entity.Review, error) {
	args := m.Called(db, productID)
	return args.Get(0).([]*entity.Review), args.Error(1)
}

func (m *ReviewRepositoryMock) GetUserReviews(db database.Database, userID uint) ([]*entity.Review, error) {
	args := m.Called(db, userID)
	return args.Get(0).([]*entity.Review), args.Error(1)
}

func (m *ReviewRepositoryMock) DeleteReviewByID(db database.Database, id uint) error {
	args := m.Called(db, id)
	return args.Error(0)
}

func (m *ReviewRepositoryMock) DeleteByProductID(db database.Database, productID uint) error {
	args := m.Called(db, productID)
	return args.Error(0)
}
