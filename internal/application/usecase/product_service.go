package usecase

import (
	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
)

type ProductService interface {
	ParseSlug(string) (string, error)
	ParseProduct(entity.Product) (productdto.ProductCredential)
	FindProductBySlug(string) (*productdto.ProductCredential, error)
	GetProduct(uint) (*productdto.ProductCredential, error)
	GetProductBySlug(string) (*productdto.ProductCredential, error)
	GetProducts() ([]productdto.ProductCredential, error)
	CreateProduct(productdto.CreateProductRequest) error
	UpdateProduct(productdto.UpdateProductRequest) error
	DeleteProduct(uint) error
	GetCategoryProducts(categoryID uint) ([]productdto.ProductCredential, error)
	UpdateProductsPrice([]productdto.ProductPriceUpdateCredentials) error
	GetProductPrices() ([]productdto.ProductPrices, error)
}
