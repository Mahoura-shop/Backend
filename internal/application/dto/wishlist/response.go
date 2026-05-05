package wishlistdto

import productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"

type WishlistItemCredential struct {
	ID        uint                       `json:"id"`
	ProductID uint                       `json:"productID"`
	Product   productdto.ProductCredential `json:"product"`
}
