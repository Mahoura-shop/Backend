package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type WishlistRepositoryMock struct {
	mock.Mock
}

func NewWishlistRepositoryMock() *WishlistRepositoryMock {
	return &WishlistRepositoryMock{}
}

func (m *WishlistRepositoryMock) AddToWishlist(db database.Database, item entity.Wishlist) error {
	args := m.Called(db, item)
	return args.Error(0)
}

func (m *WishlistRepositoryMock) RemoveFromWishlist(db database.Database, userID, productID uint) error {
	args := m.Called(db, userID, productID)
	return args.Error(0)
}

func (m *WishlistRepositoryMock) GetWishlistByUserID(db database.Database, userID uint) ([]*entity.Wishlist, error) {
	args := m.Called(db, userID)
	return args.Get(0).([]*entity.Wishlist), args.Error(1)
}

func (m *WishlistRepositoryMock) FindWishlistItem(db database.Database, userID, productID uint) (*entity.Wishlist, error) {
	args := m.Called(db, userID, productID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Wishlist), args.Error(1)
}
