package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type ProductRepository interface {
	CreateProduct(database.Database, *entity.Product) error
	FindProductByID(database.Database, uint) (*entity.Product, error)
	FindProductBySlug(database.Database, string) (*entity.Product, error)
	GetProducts(database.Database) ([]*entity.Product, error)
	DeleteProductByID(database.Database, uint) error
	UpdateProduct(database.Database, *entity.Product) error
}
