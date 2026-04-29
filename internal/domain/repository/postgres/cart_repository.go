package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type CartRepository interface {
	FindCartByUserID(database.Database, uint) (*entity.Cart, error)
	FindCartByID(database.Database, uint) (*entity.Cart, error)
	CreateCart(database.Database, entity.Cart) (error)
	DeleteCart(database.Database, uint) (error)
	DeleteCartItems(db database.Database, cartID uint) (error)
	UpdateCart(database.Database, entity.Cart) (error)
	AddProductToCart(db database.Database, productID uint, cartID uint) (error)
	RemoveProductToCart(db database.Database, productID uint, cartID uint) (error)

	FindCartItemByID(db database.Database, cartItemID uint) (*entity.CartItem, error)
	FindCartItemByProductID(db database.Database, productID uint, cartID uint) (*entity.CartItem, error)
	CreateCartItem(db database.Database, cartItem entity.CartItem) (error)
	DeleteCartItemByID(db database.Database, cartItemID uint) (error)
	UpdateCartItem(db database.Database, cartItem entity.CartItem) (error)
	IncreaseCartItemCount(db database.Database, cartItemID uint) (error)
	DecreaseCartItemCount(db database.Database, cartItemID uint) (error)
}
