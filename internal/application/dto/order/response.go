package orderdto

import productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"

type OrderCredential struct {
	ID    uint                  `json:"id"`
	itmes []OrderItemCredential `json:"items"`
}

type OrderItemCredential struct {
	ID      uint                         `json:"id"`
	Product productdto.ProductCredential `json:"product"`
	Count   uint                         `json:"count"`
	PriceSnapshot                        `json:"priceSnapshot"`
}