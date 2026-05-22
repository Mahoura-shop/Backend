package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type ReviewRepository interface {
	CreateReview(database.Database, entity.Review) (*entity.Review, error)
	FindReviewByID(database.Database, uint) (*entity.Review, error)
	FindReviewByUserAndProduct(database.Database, uint, uint) (*entity.Review, error)
	GetProductReviews(database.Database, uint) ([]*entity.Review, error)
	GetUserReviews(database.Database, uint) ([]*entity.Review, error)
	DeleteReviewByID(database.Database, uint) error
	DeleteByProductID(database.Database, uint) error
}
