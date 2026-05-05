package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type ProductImageRepository interface {
	CreateProductImage(database.Database, entity.ProductImage) (*entity.ProductImage, error)
	FindProductImageByID(database.Database, uint) (*entity.ProductImage, error)
	GetProductImages(database.Database, uint) ([]*entity.ProductImage, error)
	DeleteProductImageByID(database.Database, uint) error
}
