package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type CategoryRepository struct{}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{}
}

func (repo *CategoryRepository) FindCategoryByID(db database.Database, categoryID uint) (*entity.Category, error) {
	var category entity.Category
	result := db.GetDB().Where("id = ?", categoryID).First(&category)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &category, nil
}

func (repo *CategoryRepository) FindCategoryBySlug(db database.Database, slug string) (*entity.Category, error) {
	var category entity.Category
	result := db.GetDB().Where("slug = ?", slug).First(&category)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &category, nil
}

func (repo *CategoryRepository) GetCategories(db database.Database) ([]*entity.Category, error) {
	var categories []*entity.Category
	result := db.GetDB().Find(&categories)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return categories, nil
}

func (repo *CategoryRepository) CreateCategory(db database.Database, category *entity.Category) error {
	return db.GetDB().Create(&category).Error
}

func (repo *CategoryRepository) UpdateCategory(db database.Database, category *entity.Category) error {
	return db.GetDB().Save(&category).Error
}

func (repo *CategoryRepository) DeleteCategoryByID(db database.Database, categoryID uint) error {
	return db.GetDB().Where("id = ?", categoryID).Delete(&entity.Category{}).Error
}