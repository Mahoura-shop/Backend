package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

func (repo *ProductRepository) FindProductByID(db database.Database, productID uint) (*entity.Product, error) {
	var product entity.Product
	result := db.GetDB().Preload("Brand").Preload("Category").Where("id = ?", productID).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (repo *ProductRepository) FindProductBySlug(db database.Database, slug string) (*entity.Product, error) {
	var product entity.Product
	result := db.GetDB().Preload("Brand").Preload("Category").Where("slug = ?", slug).First(&product)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &product, nil
}

func (repo *ProductRepository) FindProductByName(db database.Database, name string) (*entity.Product, error) {
	var product entity.Product
	result := db.GetDB().Preload("Brand").Preload("Category").Where("name = ?", name).First(&product)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &product, nil
}

func (repo *ProductRepository) GetProducts(db database.Database) ([]*entity.Product, error) {
	var products []*entity.Product
	result := db.GetDB().Preload("Brand").Preload("Category").Find(&products)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return products, nil
}

func (repo *ProductRepository) CreateProduct(db database.Database, product *entity.Product) (*entity.Product, error) {
    result := db.GetDB().Create(product)
    if result.Error != nil {
        return nil, result.Error
    }
    
    return product, nil
}

func (repo *ProductRepository) UpdateProduct(db database.Database, product *entity.Product) error {
	return db.GetDB().Save(&product).Error
}

func (repo *ProductRepository) DeleteProductByID(db database.Database, productID uint) error {
	return db.GetDB().Where("id = ?", productID).Unscoped().Delete(&entity.Product{}).Error
}

func (repo *ProductRepository) GetProductsCount(db database.Database) (uint, error) {
	var count int64
	err := db.GetDB().Model(&entity.Product{}).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return uint(count), nil
}