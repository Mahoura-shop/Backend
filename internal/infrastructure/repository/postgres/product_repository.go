package postgres

import (
	domainPostgres "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
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
	result := db.GetDB().Preload("Brand").Preload("Category").Preload("Currency").Preload("Images").Where("id = ?", productID).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (repo *ProductRepository) FindProductBySlug(db database.Database, slug string) (*entity.Product, error) {
	var product entity.Product
	result := db.GetDB().Preload("Brand").Preload("Category").Preload("Currency").Preload("Images").Where("slug = ?", slug).First(&product)
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
	result := db.GetDB().Preload("Brand").Preload("Category").Preload("Currency").Where("name = ?", name).First(&product)
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
	result := db.GetDB().Preload("Brand").Preload("Category").Preload("Currency").Find(&products)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return products, nil
}


func (repo *ProductRepository) GetCategoryProducts(db database.Database, categoryID uint) ([]*entity.Product, error) {
	var products []*entity.Product
	result := db.GetDB().Preload("Brand").Preload("Category").Preload("Currency").Where("category_id = ?", categoryID).Find(&products)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return products, nil
}

func (repo *ProductRepository) CreateProduct(db database.Database, product entity.Product) (*entity.Product, error) {
    result := db.GetDB().Create(&product)
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &product, nil
}

func (repo *ProductRepository) UpdateProduct(db database.Database, product entity.Product) error {
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

func (repo *ProductRepository) SearchProducts(db database.Database, filter domainPostgres.ProductFilter) ([]*entity.Product, error) {
	var products []*entity.Product
	q := db.GetDB().Preload("Brand").Preload("Category").Preload("Currency")

	if filter.Query != "" {
		q = q.Where(
			"to_tsvector('simple', coalesce(name,'') || ' ' || coalesce(description,'')) @@ plainto_tsquery('simple', ?)",
			filter.Query,
		)
	}
	if filter.CategoryID != nil {
		q = q.Where("category_id = ?", *filter.CategoryID)
	}
	if filter.BrandID != nil {
		q = q.Where("brand_id = ?", *filter.BrandID)
	}
	if filter.MinPrice != nil {
		q = q.Where("irr_price >= ?", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		q = q.Where("irr_price <= ?", *filter.MaxPrice)
	}
	if filter.InStock != nil && *filter.InStock {
		q = q.Where("quantity > 0")
	}
	switch filter.SortBy {
	case "price_asc":
		q = q.Order("irr_price ASC")
	case "price_desc":
		q = q.Order("irr_price DESC")
	case "newest":
		q = q.Order("created_at DESC")
	default:
		q = q.Order("priority DESC, created_at DESC")
	}

	if err := q.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (repo *ProductRepository) GetRelatedProducts(db database.Database, productID, categoryID, brandID uint, limit int) ([]*entity.Product, error) {
	var products []*entity.Product
	result := db.GetDB().
		Preload("Brand").Preload("Category").Preload("Currency").
		Where("id != ? AND (category_id = ? OR brand_id = ?)", productID, categoryID, brandID).
		Where("is_active = true AND quantity > 0").
		Order("priority DESC").
		Limit(limit).
		Find(&products)
	if result.Error != nil {
		return nil, result.Error
	}
	return products, nil
}