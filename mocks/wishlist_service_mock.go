package mocks

import (
	wishlistdto "github.com/Mahoura-shop/Backend/internal/application/dto/wishlist"
	"github.com/stretchr/testify/mock"
)

type WishlistServiceMock struct {
	mock.Mock
}

func NewWishlistServiceMock() *WishlistServiceMock {
	return &WishlistServiceMock{}
}

func (m *WishlistServiceMock) AddToWishlist(userID, productID uint) error {
	args := m.Called(userID, productID)
	return args.Error(0)
}

func (m *WishlistServiceMock) RemoveFromWishlist(userID, productID uint) error {
	args := m.Called(userID, productID)
	return args.Error(0)
}

func (m *WishlistServiceMock) GetMyWishlist(userID uint) ([]wishlistdto.WishlistItemCredential, error) {
	args := m.Called(userID)
	return args.Get(0).([]wishlistdto.WishlistItemCredential), args.Error(1)
}
