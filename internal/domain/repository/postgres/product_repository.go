package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type ProductRepository interface {
	CreateProduct(database.Database, *entity.Product) error
	FindProductByID(db database.Database, productID uint) (*entity.Product, error)
	FindProductBySlug(db database.Database, slug string) (*entity.Product, error)
	GetProducts(db database.Database) ([]*entity.Product, error)
	DeleteProductByID(db database.Database, productID uint) error
	UpdateProduct(db database.Database, product *entity.Product) error
}
