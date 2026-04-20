package usecase

import (
	cartdto "github.com/Mahoura-shop/Backend/internal/application/dto/cart"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
)

type CartService interface {
	ParseCart(entity.Cart) (cartdto.CartCredential)
	GetUserCart(userID uint) (cartdto.CartCredential, error)
	AddProductToCart(cartdto.UpdateProductCountInCart) (error)
	RemoveProductFromCart(cartdto.UpdateProductCountInCart) (error)
	
	ParseCartItem(cart entity.CartItem) (cartdto.CartItemCredential)
	DeleteCartItems(cartID uint) (error)
}
