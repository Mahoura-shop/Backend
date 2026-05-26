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
	result := db.GetDB().Unscoped().Preload("Brand").Preload("Category").Preload("Currency").Preload("Images").Where("id = ?", productID).First(&product)
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

func (repo *ProductRepository) FindProductByExternalID(db database.Database, externalID string) (*entity.Product, error) {
	var product entity.Product
	result := db.GetDB().Where("external_id = ?", externalID).First(&product)
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

func (repo *ProductRepository) GetBrandProducts(db database.Database, brandID uint) ([]*entity.Product, error) {
	var products []*entity.Product
	result := db.GetDB().Preload("Brand").Preload("Category").Preload("Currency").Where("brand_id = ?", brandID).Find(&products)

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
	return db.GetDB().Omit("Brand", "Category", "Currency", "Images").Save(&product).Error
}

func (repo *ProductRepository) DeleteProductByID(db database.Database, productID uint) error {
	return db.GetDB().Where("id = ?", productID).Delete(&entity.Product{}).Error
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
	case "popularity":
		q = q.Joins("LEFT JOIN order_items ON order_items.product_id = products.id").
			Joins("LEFT JOIN orders ON orders.id = order_items.order_id AND orders.status = ?", 2).
			Group("products.id").
			Order("COALESCE(SUM(order_items.count), 0) DESC, products.priority DESC, products.created_at DESC")
	case "most_visited":
		q = q.Joins("LEFT JOIN product_visits ON product_visits.product_id = products.id").
			Group("products.id").
			Order("COALESCE(COUNT(product_visits.id), 0) DESC, products.priority DESC, products.created_at DESC")
	default:
		q = q.Order("priority DESC, created_at DESC")
	}

	if err := q.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (repo *ProductRepository) SearchProductsWithCount(db database.Database, filter domainPostgres.ProductFilter) ([]*entity.Product, int64, error) {
	var products []*entity.Product
	var count int64
	q := db.GetDB().Preload("Brand").Preload("Category").Preload("Currency")

	q = q.Where("is_active = true")

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

	if err := q.Model(&entity.Product{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	switch filter.SortBy {
	case "price_asc":
		q = q.Order("irr_price ASC")
	case "price_desc":
		q = q.Order("irr_price DESC")
	case "newest":
		q = q.Order("created_at DESC")
	case "popularity":
		q = q.Joins("LEFT JOIN order_items ON order_items.product_id = products.id").
			Joins("LEFT JOIN orders ON orders.id = order_items.order_id AND orders.status = ?", 2).
			Group("products.id").
			Order("COALESCE(SUM(order_items.count), 0) DESC, products.priority DESC, products.created_at DESC")
	case "most_visited":
		q = q.Joins("LEFT JOIN product_visits ON product_visits.product_id = products.id").
			Group("products.id").
			Order("COALESCE(COUNT(product_visits.id), 0) DESC, products.priority DESC, products.created_at DESC")
	default:
		q = q.Order("priority DESC, created_at DESC")
	}

	if filter.Limit != nil {
		q = q.Limit(*filter.Limit)
	}
	if filter.Offset != nil {
		q = q.Offset(*filter.Offset)
	}

	if err := q.Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, count, nil
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

func (repo *ProductRepository) GetLowStockProducts(db database.Database, limit int) ([]domainPostgres.LowStockProduct, error) {
	var results []domainPostgres.LowStockProduct
	result := db.GetDB().
		Model(&entity.Product{}).
		Select("id, name, quantity, min_order").
		Where("quantity < min_order AND is_active = true").
		Order("quantity ASC").
		Limit(limit).
		Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}
	return results, nil
}

func (repo *ProductRepository) GetTopSoldProducts(db database.Database, limit int) ([]domainPostgres.TopSoldProduct, error) {
	var results []domainPostgres.TopSoldProduct
	result := db.GetDB().
		Model(&entity.OrderItem{}).
		Select("product_id as id, products.name, SUM(order_items.count) as quantity, SUM(order_items.price_snapshot * order_items.count) as revenue").
		Joins("LEFT JOIN products ON products.id = order_items.product_id").
		Group("order_items.product_id, products.name").
		Order("quantity DESC").
		Limit(limit).
		Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}
	return results, nil
}