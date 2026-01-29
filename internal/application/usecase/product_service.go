package usecase

import (
	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
)

type ProductService interface {
	ParseSlug(string) (string, error)
	FindProductBySlug(string) (*entity.Product, error)
	GetProducts() ([]productdto.ProductCredential, error)
	CreateProduct(productdto.CreateProductRequest) error
	UpdateProduct(productdto.UpdateProductRequest) error
	DeleteProduct(uint) error
}
