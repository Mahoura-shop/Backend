package mocks

import (
	cartdto "github.com/Mahoura-shop/Backend/internal/application/dto/cart"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type CartServiceMock struct {
	mock.Mock
}

func NewCartServiceMock() *CartServiceMock {
	return &CartServiceMock{}
}

func (m *CartServiceMock) ParseCart(c entity.Cart) cartdto.CartCredential {
	args := m.Called(c)
	return args.Get(0).(cartdto.CartCredential)
}

func (m *CartServiceMock) GetUserCart(userID uint) (cartdto.CartCredential, error) {
	args := m.Called(userID)
	return args.Get(0).(cartdto.CartCredential), args.Error(1)
}

func (m *CartServiceMock) AddProductToCart(req cartdto.UpdateProductCountInCart) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *CartServiceMock) RemoveProductFromCart(req cartdto.UpdateProductCountInCart) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *CartServiceMock) ParseCartItem(item entity.CartItem) cartdto.CartItemCredential {
	args := m.Called(item)
	return args.Get(0).(cartdto.CartItemCredential)
}

func (m *CartServiceMock) DeleteCartItems(cartID uint) error {
	args := m.Called(cartID)
	return args.Error(0)
}
