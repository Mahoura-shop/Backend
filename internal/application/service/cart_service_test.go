package service

import (
	"errors"
	"testing"

	cartdto "github.com/Mahoura-shop/Backend/internal/application/dto/cart"
	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"
	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/suite"
)

type CartServiceTestSuite struct {
	suite.Suite
	constants      *bootstrap.Constants
	cartRepo       *mocks.CartRepositoryMock
	productRepo    *mocks.ProductRepositoryMock
	userService    *mocks.UserServiceMock
	productService *mocks.ProductServiceMock
	db             *mocks.DatabaseMock
	service        *CartService
}

func (s *CartServiceTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.cartRepo = mocks.NewCartRepositoryMock()
	s.productRepo = mocks.NewProductRepositoryMock()
	s.userService = mocks.NewUserServiceMock()
	s.productService = mocks.NewProductServiceMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewCartService(CartServiceDeps{
		Constants:         s.constants,
		CartRepository:    s.cartRepo,
		ProductRepository: s.productRepo,
		UserService:       s.userService,
		ProductService:    s.productService,
		DB:                s.db,
	})
}

func (s *CartServiceTestSuite) TestGetUserCart_NoCart_ReturnsEmpty() {
	var nilCart *entity.Cart
	s.cartRepo.On("FindCartByUserID", s.db, uint(1)).Return(nilCart, nil).Once()

	result, err := s.service.GetUserCart(1)

	s.NoError(err)
	s.Empty(result.Items)
	s.cartRepo.AssertExpectations(s.T())
}

func (s *CartServiceTestSuite) TestGetUserCart_RepoError() {
	s.cartRepo.On("FindCartByUserID", s.db, uint(1)).Return(nil, errors.New("db error")).Once()

	_, err := s.service.GetUserCart(1)

	s.Error(err)
}

func (s *CartServiceTestSuite) TestGetUserCart_ExistingCart_ReturnsItems() {
	cart := &entity.Cart{
		Model:  dbmodel.Model{ID: 3},
		UserID: 1,
		Items:  []entity.CartItem{},
	}
	s.cartRepo.On("FindCartByUserID", s.db, uint(1)).Return(cart, nil).Once()
	s.userService.On("ParseUser", cart.User).Return(userdto.UserCredential{}).Once()

	result, err := s.service.GetUserCart(1)

	s.NoError(err)
	s.Equal(uint(3), result.ID)
}

func (s *CartServiceTestSuite) TestAddProductToCart_ProductNotFound() {
	var nilProduct *entity.Product
	s.productRepo.On("FindProductByID", s.db, uint(5)).Return(nilProduct, nil).Once()

	err := s.service.AddProductToCart(cartdto.UpdateProductCountInCart{
		UserID:    1,
		ProductID: 5,
	})

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.cartRepo.AssertNotCalled(s.T(), "AddProductToCart")
}

func (s *CartServiceTestSuite) TestAddProductToCart_ExceedsStock_ValidationError() {
	product := &entity.Product{
		Model:    dbmodel.Model{ID: 5},
		Quantity: 2,
	}
	cart := &entity.Cart{Model: dbmodel.Model{ID: 10}}
	existingItem := &entity.CartItem{ProductID: 5, CartID: 10, Count: 2}

	s.productRepo.On("FindProductByID", s.db, uint(5)).Return(product, nil).Once()
	s.cartRepo.On("FindCartByUserID", s.db, uint(1)).Return(cart, nil).Once()
	s.cartRepo.On("FindCartItemByProductID", s.db, uint(5), uint(10)).Return(existingItem, nil).Once()

	err := s.service.AddProductToCart(cartdto.UpdateProductCountInCart{
		UserID:    1,
		ProductID: 5,
	})

	s.Error(err)
	var ve exception.ValidationErrors
	s.True(errors.As(err, &ve))
}

func (s *CartServiceTestSuite) TestAddProductToCart_NewCart_CreatesCartThenAdds() {
	product := &entity.Product{
		Model:    dbmodel.Model{ID: 5},
		Quantity: 10,
	}
	var nilCart *entity.Cart
	newCart := &entity.Cart{Model: dbmodel.Model{ID: 20}}
	var nilItem *entity.CartItem

	s.productRepo.On("FindProductByID", s.db, uint(5)).Return(product, nil).Once()
	s.cartRepo.On("FindCartByUserID", s.db, uint(1)).Return(nilCart, nil).Once()
	s.cartRepo.On("CreateCart", s.db, entity.Cart{UserID: 1}).Return(nil).Once()
	s.cartRepo.On("FindCartByUserID", s.db, uint(1)).Return(newCart, nil).Once()
	s.cartRepo.On("FindCartItemByProductID", s.db, uint(5), uint(20)).Return(nilItem, nil).Once()
	s.cartRepo.On("AddProductToCart", s.db, uint(5), uint(20)).Return(nil).Once()

	err := s.service.AddProductToCart(cartdto.UpdateProductCountInCart{
		UserID:    1,
		ProductID: 5,
	})

	s.NoError(err)
	s.cartRepo.AssertExpectations(s.T())
}

func (s *CartServiceTestSuite) TestRemoveProductFromCart_ProductNotFound() {
	var nilProduct *entity.Product
	s.productRepo.On("FindProductByID", s.db, uint(5)).Return(nilProduct, nil).Once()

	err := s.service.RemoveProductFromCart(cartdto.UpdateProductCountInCart{
		UserID:    1,
		ProductID: 5,
	})

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
}

func (s *CartServiceTestSuite) TestRemoveProductFromCart_NoCart_NotFound() {
	product := &entity.Product{Model: dbmodel.Model{ID: 5}}
	var nilCart *entity.Cart

	s.productRepo.On("FindProductByID", s.db, uint(5)).Return(product, nil).Once()
	s.cartRepo.On("FindCartByUserID", s.db, uint(1)).Return(nilCart, nil).Once()

	err := s.service.RemoveProductFromCart(cartdto.UpdateProductCountInCart{
		UserID:    1,
		ProductID: 5,
	})

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
}

func (s *CartServiceTestSuite) TestRemoveProductFromCart_LastItem_DeletesItem() {
	product := &entity.Product{Model: dbmodel.Model{ID: 5}}
	cart := &entity.Cart{Model: dbmodel.Model{ID: 10}}
	item := &entity.CartItem{Model: dbmodel.Model{ID: 7}, ProductID: 5, CartID: 10, Count: 1}

	s.productRepo.On("FindProductByID", s.db, uint(5)).Return(product, nil).Once()
	s.cartRepo.On("FindCartByUserID", s.db, uint(1)).Return(cart, nil).Once()
	s.cartRepo.On("FindCartItemByProductID", s.db, uint(5), uint(10)).Return(item, nil).Once()
	s.cartRepo.On("DeleteCartItemByID", s.db, uint(7)).Return(nil).Once()

	err := s.service.RemoveProductFromCart(cartdto.UpdateProductCountInCart{
		UserID:    1,
		ProductID: 5,
	})

	s.NoError(err)
	s.cartRepo.AssertExpectations(s.T())
}

func (s *CartServiceTestSuite) TestRemoveProductFromCart_MultipleItems_Decreases() {
	product := &entity.Product{Model: dbmodel.Model{ID: 5}}
	cart := &entity.Cart{Model: dbmodel.Model{ID: 10}}
	item := &entity.CartItem{Model: dbmodel.Model{ID: 7}, Count: 3}

	s.productRepo.On("FindProductByID", s.db, uint(5)).Return(product, nil).Once()
	s.cartRepo.On("FindCartByUserID", s.db, uint(1)).Return(cart, nil).Once()
	s.cartRepo.On("FindCartItemByProductID", s.db, uint(5), uint(10)).Return(item, nil).Once()
	s.cartRepo.On("DecreaseCartItemCount", s.db, uint(7)).Return(nil).Once()

	err := s.service.RemoveProductFromCart(cartdto.UpdateProductCountInCart{
		UserID:    1,
		ProductID: 5,
	})

	s.NoError(err)
	s.cartRepo.AssertExpectations(s.T())
}

func (s *CartServiceTestSuite) TestDeleteCartItems_CallsRepo() {
	s.cartRepo.On("DeleteCartItems", s.db, uint(10)).Return(nil).Once()

	err := s.service.DeleteCartItems(10)

	s.NoError(err)
	s.cartRepo.AssertExpectations(s.T())
}

func (s *CartServiceTestSuite) TestParseCartItem_MapsFields() {
	cartItem := entity.CartItem{
		Model:   dbmodel.Model{ID: 3},
		Count:   2,
		Product: entity.Product{},
	}
	parsedProduct := productdto.ProductCredential{ID: 5, Name: "Cream"}
	s.productService.On("ParseProduct", cartItem.Product).Return(parsedProduct).Once()

	result := s.service.ParseCartItem(cartItem)

	s.Equal(uint(3), result.ID)
	s.Equal(uint(2), result.Count)
	s.Equal(parsedProduct, result.Product)
	s.productService.AssertExpectations(s.T())
}

func TestCartServiceSuite(t *testing.T) {
	suite.Run(t, new(CartServiceTestSuite))
}
