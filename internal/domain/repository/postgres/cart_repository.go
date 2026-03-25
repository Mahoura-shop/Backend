package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type CartRepository interface {
	FindCartByUserID(database.Database, uint) (*entity.Cart, error)
	FindCartByID(database.Database, uint) (*entity.Cart, error)
	CreateCart(database.Database, entity.Cart) (error)
	DeleteCartByID(database.Database, uint) (error)
	UpdateCart(database.Database, entity.Cart) (error)
	AddProductToCart(db database.Database, productID uint, cartID uint) (error)
}
