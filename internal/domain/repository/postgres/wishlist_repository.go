package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type WishlistRepository interface {
	AddToWishlist(database.Database, entity.Wishlist) error
	RemoveFromWishlist(database.Database, uint, uint) error
	GetWishlistByUserID(database.Database, uint) ([]*entity.Wishlist, error)
	FindWishlistItem(database.Database, uint, uint) (*entity.Wishlist, error)
	GetWishlistByProductID(database.Database, uint) ([]*entity.Wishlist, error)
}
