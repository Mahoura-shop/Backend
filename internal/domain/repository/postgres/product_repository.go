package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type ProductFilter struct {
	Query      string
	CategoryID *uint
	BrandID    *uint
	MinPrice   *uint
	MaxPrice   *uint
	InStock    *bool
	SortBy     string
	Limit      *int
	Offset     *int
}

type LowStockProduct struct {
	ID       uint
	Name     string
	Quantity uint
	MinOrder uint
}

type TopSoldProduct struct {
	ID       uint
	Name     string
	Quantity uint
	Revenue  uint
}

type ProductRepository interface {
	CreateProduct(database.Database, entity.Product) (*entity.Product, error)
	FindProductByID(database.Database, uint) (*entity.Product, error)
	FindProductBySlug(database.Database, string) (*entity.Product, error)
	FindProductByName(database.Database, string) (*entity.Product, error)
	GetProducts(database.Database) ([]*entity.Product, error)
	SearchProducts(database.Database, ProductFilter) ([]*entity.Product, error)
	SearchProductsWithCount(database.Database, ProductFilter) ([]*entity.Product, int64, error)
	GetRelatedProducts(database.Database, uint, uint, uint, int) ([]*entity.Product, error)
	GetCategoryProducts(database.Database, uint) ([]*entity.Product, error)
	GetProductsCount(database.Database) (uint, error)
	DeleteProductByID(database.Database, uint) error
	UpdateProduct(database.Database, entity.Product) error
	GetLowStockProducts(database.Database, int) ([]LowStockProduct, error)
	GetTopSoldProducts(database.Database, int) ([]TopSoldProduct, error)
}
