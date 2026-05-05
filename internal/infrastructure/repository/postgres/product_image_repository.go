package postgres

import (
	"errors"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type ProductImageRepository struct{}

func NewProductImageRepository() *ProductImageRepository {
	return &ProductImageRepository{}
}

func (r *ProductImageRepository) CreateProductImage(db database.Database, img entity.ProductImage) (*entity.ProductImage, error) {
	result := db.GetDB().Create(&img)
	if result.Error != nil {
		return nil, result.Error
	}
	return &img, nil
}

func (r *ProductImageRepository) FindProductImageByID(db database.Database, imageID uint) (*entity.ProductImage, error) {
	var img entity.ProductImage
	result := db.GetDB().Where("id = ?", imageID).First(&img)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &img, nil
}

func (r *ProductImageRepository) GetProductImages(db database.Database, productID uint) ([]*entity.ProductImage, error) {
	var images []*entity.ProductImage
	result := db.GetDB().Where("product_id = ?", productID).Order("position ASC, created_at ASC").Find(&images)
	if result.Error != nil {
		return nil, result.Error
	}
	return images, nil
}

func (r *ProductImageRepository) DeleteProductImageByID(db database.Database, imageID uint) error {
	return db.GetDB().Where("id = ?", imageID).Delete(&entity.ProductImage{}).Error
}
