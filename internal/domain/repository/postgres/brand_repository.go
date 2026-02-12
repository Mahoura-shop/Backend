package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type BrandRepository interface {
	GetBrandProductsCount(db database.Database, brandID uint) (uint, error)
	CreateBrand(database.Database, *entity.Brand) (*entity.Brand, error)
	FindBrandByID(database.Database, uint) (*entity.Brand, error)
	FindBrandBySlug(database.Database, string) (*entity.Brand, error)
	GetBrands(database.Database) ([]*entity.Brand, error)
	DeleteBrandByID(database.Database, uint) error
	UpdateBrand(database.Database, *entity.Brand) error
}
