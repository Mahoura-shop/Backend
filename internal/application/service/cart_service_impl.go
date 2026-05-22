package service

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	cartdto "github.com/Mahoura-shop/Backend/internal/application/dto/cart"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/domain/s3"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type CartService struct {
	constants          *bootstrap.Constants
	cartRepository     postgres.CartRepository
	productRepository  postgres.ProductRepository
	userService        usecase.UserService
	productService     usecase.ProductService
	s3Storage          s3.S3Storage
	db                 database.Database
}

type CartServiceDeps struct {
	Constants          *bootstrap.Constants
	CartRepository     postgres.CartRepository
	ProductRepository  postgres.ProductRepository
	UserService        usecase.UserService
	ProductService     usecase.ProductService
	S3Storage          s3.S3Storage
	DB                 database.Database
}

func NewCartService(deps CartServiceDeps) *CartService {
	return &CartService{
		constants:          deps.Constants,
		cartRepository:     deps.CartRepository,
		productRepository:  deps.ProductRepository,
		userService:        deps.UserService,
		productService:     deps.ProductService,
		s3Storage:          deps.S3Storage,
		db:                 deps.DB,
	}
}

func (cartService *CartService) ParseCart(cart entity.Cart) (cartdto.CartCredential) {
	response := cartdto.CartCredential{
		ID:          cart.ID,
	}
	
	response.User = cartService.userService.ParseUser(cart.User)

	responses := []cartdto.CartItemCredential{}
	for _, item := range cart.Items {
		parsedItem := cartService.ParseCartItem(item)
		responses = append(responses, parsedItem)
	}
	response.Items = responses

	return response
}

func (cartService *CartService) GetUserCart(userID uint) (cartdto.CartCredential, error) {
	cart, err := cartService.cartRepository.FindCartByUserID(cartService.db, userID)
	if err != nil {
		return cartdto.CartCredential{}, err
	}

	if cart == nil {
		return cartdto.CartCredential{Items: []cartdto.CartItemCredential{}}, nil
	}

	return cartService.ParseCart(*cart), nil
}

func (cartService *CartService) AddProductToCart(addProductToCartInfo cartdto.UpdateProductCountInCart) (error) {
	product, err := cartService.productRepository.FindProductByID(cartService.db, addProductToCartInfo.ProductID)
	if err != nil {
		return err
	}
	if product == nil {
		return exception.NotFoundError{Item: cartService.constants.Field.Product}
	}

	cart, err := cartService.cartRepository.FindCartByUserID(cartService.db, addProductToCartInfo.UserID)
	if err != nil {
		return err
	}
	if cart == nil {
		newCart := entity.Cart{UserID: addProductToCartInfo.UserID}
		if err = cartService.cartRepository.CreateCart(cartService.db, newCart); err != nil {
			return err
		}
		cart, err = cartService.cartRepository.FindCartByUserID(cartService.db, addProductToCartInfo.UserID)
		if err != nil {
			return err
		}
	}

	currentCartItem, err := cartService.cartRepository.FindCartItemByProductID(cartService.db, product.ID, cart.ID)
	if err != nil {
		return err
	}

	currentCount := uint(0)
	if currentCartItem != nil {
		currentCount = currentCartItem.Count
	}

	if currentCount+1 > product.Quantity {
		var validationErrors exception.ValidationErrors
		validationErrors.Add(cartService.constants.Field.CartItem, cartService.constants.Tag.Invalid)
		return validationErrors
	}

	return cartService.cartRepository.AddProductToCart(cartService.db, product.ID, cart.ID)
}

func (cartService *CartService) RemoveProductFromCart(removeProductFromCartInfo cartdto.UpdateProductCountInCart) (error) {
	product, err := cartService.productRepository.FindProductByID(cartService.db, removeProductFromCartInfo.ProductID)
	if err != nil {
		return err
	}
	if product == nil {
		return exception.NotFoundError{Item: cartService.constants.Field.Product}
	}

	cart, err := cartService.cartRepository.FindCartByUserID(cartService.db, removeProductFromCartInfo.UserID)
	if err != nil {
		return err
	}
	if cart == nil {
		return exception.NotFoundError{Item: cartService.constants.Field.CartItem}
	}

	cartItem, err := cartService.cartRepository.FindCartItemByProductID(cartService.db, product.ID, cart.ID)
	if err != nil {
		return err
	}
	if cartItem == nil {
		return exception.NotFoundError{Item: cartService.constants.Field.CartItem}
	}
	if cartItem.Count == 1 {
		return cartService.cartRepository.DeleteCartItemByID(cartService.db, cartItem.ID)
	}
	return cartService.cartRepository.DecreaseCartItemCount(cartService.db, cartItem.ID)
}

func (cartService *CartService) ParseCartItem(cartItem entity.CartItem) (cartdto.CartItemCredential) {
	response := cartdto.CartItemCredential{
		ID:    cartItem.ID,
		Count: cartItem.Count,
	}

	response.Product = cartService.productService.ParseProduct(cartItem.Product)

	return response
}

func (cartService *CartService) DeleteCartItems(cartID uint) (error) {
	return cartService.cartRepository.DeleteCartItems(cartService.db, cartID)
}