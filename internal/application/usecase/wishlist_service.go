package usecase

import wishlistdto "github.com/Mahoura-shop/Backend/internal/application/dto/wishlist"

type WishlistService interface {
	AddToWishlist(userID, productID uint) error
	RemoveFromWishlist(userID, productID uint) error
	GetMyWishlist(userID uint) ([]wishlistdto.WishlistItemCredential, error)
}
