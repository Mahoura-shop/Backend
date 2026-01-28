package usecase

import productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"

type ProductService interface {
	CreateProduct(productdto.CreateProductRequest) error
	GetProducts() ([]productdto.ProductCredential, error)
	DeleteProduct(uint) error
}
