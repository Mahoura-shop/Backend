package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type CartItemRepository interface {
	FindCartItemByID(database.Database, uint) (*entity.CartItem, error)
	FindCartItemByProductID(database.Database, uint) (*entity.CartItem, error)
	CreateCartItem(database.Database, entity.CartItem) (error)
	DeleteCartItemByID(database.Database, uint) (error)
	UpdateCartItem(database.Database, entity.CartItem) (error)
	IncreaseCartItemCount(database.Database, uint) (error)
	DecreaseCartItemCount(database.Database, uint) (error)
}
