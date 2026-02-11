package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type BrandRepository struct{}

func NewBrandRepository() *BrandRepository {
	return &BrandRepository{}
}

func (repo *BrandRepository) FindBrandByID(db database.Database, brandID uint) (*entity.Brand, error) {
	var brand entity.Brand
	result := db.GetDB().Where("id = ?", brandID).First(&brand)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &brand, nil
}

func (repo *BrandRepository) FindBrandBySlug(db database.Database, slug string) (*entity.Brand, error) {
	var brand entity.Brand
	result := db.GetDB().Where("slug = ?", slug).First(&brand)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &brand, nil
}

func (repo *BrandRepository) GetBrands(db database.Database) ([]*entity.Brand, error) {
	var brands []*entity.Brand
	result := db.GetDB().Find(&brands)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return brands, nil
}

func (repo *BrandRepository) CreateBrand(db database.Database, brand *entity.Brand) (*entity.Brand, error) {
	result := db.GetDB().Create(&brand)
	if result.Error != nil {
        return nil, result.Error
    }
    
    return brand, nil
}

func (repo *BrandRepository) UpdateBrand(db database.Database, brand *entity.Brand) error {
	return db.GetDB().Save(&brand).Error
}

func (repo *BrandRepository) DeleteBrandByID(db database.Database, brandID uint) error {
	return db.GetDB().Where("id = ?", brandID).Unscoped().Delete(&entity.Brand{}).Error
}