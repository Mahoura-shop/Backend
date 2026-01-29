package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type BrandRepository interface {
	CreateBrand(db database.Database, brand *entity.Brand) error
	FindBrandByID(db database.Database, brandID uint) (*entity.Brand, error)
	FindBrandBySlug(db database.Database, slug string) (*entity.Brand, error)
	GetBrands(db database.Database) ([]*entity.Brand, error)
	DeleteBrandByID(db database.Database, brandID uint) error
	UpdateBrand(db database.Database, brand *entity.Brand) error
}
