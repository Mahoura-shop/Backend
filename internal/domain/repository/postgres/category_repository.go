package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type CategoryRepository interface {
	CreateCategory(db database.Database, category *entity.Category) error
	FindCategoryByID(db database.Database, categoryID uint) (*entity.Category, error)
	FindCategoryBySlug(db database.Database, slug string) (*entity.Category, error)
	GetCategories(db database.Database) ([]*entity.Category, error)
	DeleteCategoryByID(db database.Database, categoryID uint) error
	UpdateCategory(db database.Database, category *entity.Category) error
}
