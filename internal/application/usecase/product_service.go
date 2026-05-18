package usecase

import (
	"mime/multipart"

	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
)

type ProductService interface {
	ParseSlug(string) (string, error)
	ParseProduct(entity.Product) productdto.ProductCredential
	FindProductBySlug(string) (*productdto.ProductCredential, error)
	GetProduct(uint) (*productdto.ProductCredential, error)
	GetProductBySlug(string) (*productdto.ProductCredential, error)
	GetProducts() ([]productdto.ProductCredential, error)
	SearchProducts(productdto.ProductFilterRequest) ([]productdto.ProductCredential, error)
	SearchProductsWithPagination(productdto.ProductFilterRequest) (*productdto.ProductSearchResponse, error)
	GetRelatedProducts(uint, int) ([]productdto.ProductCredential, error)
	CreateProduct(productdto.CreateProductRequest) error
	UpdateProduct(productdto.UpdateProductRequest) error
	DeleteProduct(uint) error
	GetCategoryProducts(uint) ([]productdto.ProductCredential, error)
	UpdateProductsPrice([]productdto.ProductPriceUpdateCredentials) error
	UpdateProductsStock([]productdto.ProductStockUpdateCredentials, string) error
	GetProductPrices() ([]productdto.ProductPrices, error)
	AddProductImage(uint, *multipart.FileHeader) error
	DeleteProductImage(uint, uint) error
}
