package postgres

import (
	"time"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	postgresrepo "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type ProductVisitRepository struct{}

func NewProductVisitRepository() *ProductVisitRepository {
	return &ProductVisitRepository{}
}

func (repo *ProductVisitRepository) CreateVisit(db database.Database, visit entity.ProductVisit) (*entity.ProductVisit, error) {
	result := db.GetDB().Create(&visit)
	if result.Error != nil {
		return nil, result.Error
	}
	return &visit, nil
}

func (repo *ProductVisitRepository) GetVisitCountByProductID(db database.Database, productID uint) (int64, error) {
	var count int64
	err := db.GetDB().Model(&entity.ProductVisit{}).Where("product_id = ?", productID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *ProductVisitRepository) HasVisitedInLast24h(db database.Database, productID uint, visitorIP string) (bool, error) {
	var count int64
	oneDayAgo := time.Now().AddDate(0, 0, -1).Unix()
	err := db.GetDB().Model(&entity.ProductVisit{}).
		Where("product_id = ? AND visitor_ip = ? AND visited_at > ?", productID, visitorIP, oneDayAgo).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (repo *ProductVisitRepository) GetVisitsPerDay(db database.Database, productID uint, days int) ([]postgresrepo.VisitsByDay, error) {
	var results []postgresrepo.VisitsByDay
	result := db.GetDB().
		Model(&entity.ProductVisit{}).
		Select("DATE(TO_TIMESTAMP(visited_at)) as date, COUNT(*) as count").
		Where("product_id = ? AND visited_at >= EXTRACT(EPOCH FROM NOW() - ? * INTERVAL '1 day')", productID, days).
		Group("DATE(TO_TIMESTAMP(visited_at))").
		Order("DATE(TO_TIMESTAMP(visited_at))").
		Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}
	return results, nil
}

func (repo *ProductVisitRepository) GetAllVisitsPerDay(db database.Database, days int) ([]postgresrepo.VisitsByDay, error) {
	var results []postgresrepo.VisitsByDay
	result := db.GetDB().
		Model(&entity.ProductVisit{}).
		Select("DATE(TO_TIMESTAMP(visited_at)) as date, COUNT(*) as count").
		Where("visited_at >= EXTRACT(EPOCH FROM NOW() - ? * INTERVAL '1 day')", days).
		Group("DATE(TO_TIMESTAMP(visited_at))").
		Order("DATE(TO_TIMESTAMP(visited_at))").
		Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}
	return results, nil
}

func (repo *ProductVisitRepository) GetCategoryVisitsPerDay(db database.Database, categoryID uint, days int) ([]postgresrepo.VisitsByDay, error) {
	var results []postgresrepo.VisitsByDay
	result := db.GetDB().
		Model(&entity.ProductVisit{}).
		Select("DATE(TO_TIMESTAMP(product_visits.visited_at)) as date, COUNT(*) as count").
		Joins("JOIN products ON products.id = product_visits.product_id").
		Where("products.category_id = ? AND product_visits.visited_at >= EXTRACT(EPOCH FROM NOW() - ? * INTERVAL '1 day')", categoryID, days).
		Group("DATE(TO_TIMESTAMP(product_visits.visited_at))").
		Order("DATE(TO_TIMESTAMP(product_visits.visited_at))").
		Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}
	return results, nil
}

func (repo *ProductVisitRepository) GetBrandVisitsPerDay(db database.Database, brandID uint, days int) ([]postgresrepo.VisitsByDay, error) {
	var results []postgresrepo.VisitsByDay
	result := db.GetDB().
		Model(&entity.ProductVisit{}).
		Select("DATE(TO_TIMESTAMP(product_visits.visited_at)) as date, COUNT(*) as count").
		Joins("JOIN products ON products.id = product_visits.product_id").
		Where("products.brand_id = ? AND product_visits.visited_at >= EXTRACT(EPOCH FROM NOW() - ? * INTERVAL '1 day')", brandID, days).
		Group("DATE(TO_TIMESTAMP(product_visits.visited_at))").
		Order("DATE(TO_TIMESTAMP(product_visits.visited_at))").
		Scan(&results)
	if result.Error != nil {
		return nil, result.Error
	}
	return results, nil
}
