package postgres

import (
	"errors"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type ReviewRepository struct{}

func NewReviewRepository() *ReviewRepository {
	return &ReviewRepository{}
}

func (r *ReviewRepository) CreateReview(db database.Database, review entity.Review) (*entity.Review, error) {
	result := db.GetDB().Create(&review)
	if result.Error != nil {
		return nil, result.Error
	}
	return &review, nil
}

func (r *ReviewRepository) FindReviewByID(db database.Database, reviewID uint) (*entity.Review, error) {
	var review entity.Review
	result := db.GetDB().Preload("User").Where("id = ?", reviewID).First(&review)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &review, nil
}

func (r *ReviewRepository) FindReviewByUserAndProduct(db database.Database, userID, productID uint) (*entity.Review, error) {
	var review entity.Review
	result := db.GetDB().Where("user_id = ? AND product_id = ?", userID, productID).First(&review)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &review, nil
}

func (r *ReviewRepository) GetProductReviews(db database.Database, productID uint) ([]*entity.Review, error) {
	var reviews []*entity.Review
	result := db.GetDB().Preload("User").Where("product_id = ?", productID).Order("created_at DESC").Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}
	return reviews, nil
}

func (r *ReviewRepository) DeleteReviewByID(db database.Database, reviewID uint) error {
	return db.GetDB().Where("id = ?", reviewID).Delete(&entity.Review{}).Error
}
