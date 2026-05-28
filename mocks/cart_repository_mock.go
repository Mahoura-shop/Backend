package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type CartRepositoryMock struct {
	mock.Mock
}

func NewCartRepositoryMock() *CartRepositoryMock {
	return &CartRepositoryMock{}
}

func (m *CartRepositoryMock) FindCartByUserID(db database.Database, userID uint) (*entity.Cart, error) {
	args := m.Called(db, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Cart), args.Error(1)
}

func (m *CartRepositoryMock) FindCartByID(db database.Database, id uint) (*entity.Cart, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Cart), args.Error(1)
}

func (m *CartRepositoryMock) CreateCart(db database.Database, cart entity.Cart) error {
	args := m.Called(db, cart)
	return args.Error(0)
}

func (m *CartRepositoryMock) DeleteCart(db database.Database, cartID uint) error {
	args := m.Called(db, cartID)
	return args.Error(0)
}

func (m *CartRepositoryMock) DeleteCartItems(db database.Database, cartID uint) error {
	args := m.Called(db, cartID)
	return args.Error(0)
}

func (m *CartRepositoryMock) UpdateCart(db database.Database, cart entity.Cart) error {
	args := m.Called(db, cart)
	return args.Error(0)
}

func (m *CartRepositoryMock) AddProductToCart(db database.Database, productID uint, cartID uint) error {
	args := m.Called(db, productID, cartID)
	return args.Error(0)
}

func (m *CartRepositoryMock) RemoveProductToCart(db database.Database, productID uint, cartID uint) error {
	args := m.Called(db, productID, cartID)
	return args.Error(0)
}

func (m *CartRepositoryMock) FindCartItemByID(db database.Database, cartItemID uint) (*entity.CartItem, error) {
	args := m.Called(db, cartItemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CartItem), args.Error(1)
}

func (m *CartRepositoryMock) FindCartItemByProductID(db database.Database, productID uint, cartID uint) (*entity.CartItem, error) {
	args := m.Called(db, productID, cartID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.CartItem), args.Error(1)
}

func (m *CartRepositoryMock) CreateCartItem(db database.Database, cartItem entity.CartItem) error {
	args := m.Called(db, cartItem)
	return args.Error(0)
}

func (m *CartRepositoryMock) DeleteCartItemByID(db database.Database, cartItemID uint) error {
	args := m.Called(db, cartItemID)
	return args.Error(0)
}

func (m *CartRepositoryMock) DeleteCartItemsByProductID(db database.Database, productID uint) error {
	args := m.Called(db, productID)
	return args.Error(0)
}

func (m *CartRepositoryMock) UpdateCartItem(db database.Database, cartItem entity.CartItem) error {
	args := m.Called(db, cartItem)
	return args.Error(0)
}

func (m *CartRepositoryMock) IncreaseCartItemCount(db database.Database, cartItemID uint) error {
	args := m.Called(db, cartItemID)
	return args.Error(0)
}

func (m *CartRepositoryMock) DecreaseCartItemCount(db database.Database, cartItemID uint) error {
	args := m.Called(db, cartItemID)
	return args.Error(0)
}
