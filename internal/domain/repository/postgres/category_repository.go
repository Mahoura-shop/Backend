package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type CategoryRepository interface {
	FindCategoryByID(database.Database, uint) (*entity.Category, error)
	FindCategoryBySlug(database.Database, string) (*entity.Category, error)
	CreateCategory(database.Database, *entity.Category) error
	GetCategories(database.Database) ([]*entity.Category, error)
	DeleteCategoryByID(database.Database, uint) error
	UpdateCategory(database.Database, *entity.Category) error
}
