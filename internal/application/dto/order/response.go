package orderdto

import productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"

type OrderCredential struct {
	ID    uint                  `json:"id"`
	Itmes []OrderItemCredential `json:"items"`
}

type OrderItemCredential struct {
	ID            uint                         `json:"id"`
	Product       productdto.ProductCredential `json:"product"`
	Count         uint                         `json:"count"`
	PriceSnapshot uint                         `json:"priceSnapshot"`
}