package cartdto

import (
	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"
)

type CartCredential struct {
	ID    uint                   `json:"id"`
	User  userdto.UserCredential `json:"user"`
	Items []CartItemCredential   `json:"items"`
}

type CartItemCredential struct {
	ID uint                              `json:"id"`
	Product productdto.ProductCredential `json:"product"`
	Count uint                           `json:"count"`
}