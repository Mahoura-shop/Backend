package service

import (
	reviewdto "github.com/Mahoura-shop/Backend/internal/application/dto/review"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	domainPostgres "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type ReviewService struct {
	reviewRepository domainPostgres.ReviewRepository
	db               database.Database
}

type ReviewServiceDeps struct {
	ReviewRepository domainPostgres.ReviewRepository
	DB               database.Database
}

func NewReviewService(deps ReviewServiceDeps) *ReviewService {
	return &ReviewService{
		reviewRepository: deps.ReviewRepository,
		db:               deps.DB,
	}
}

func (s *ReviewService) isVerifiedBuyer(userID, productID uint) bool {
	var count int64
	s.db.GetDB().
		Table("order_items").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.user_id = ? AND order_items.product_id = ? AND orders.status = ?", userID, productID, enum.OrderStatusDelivered).
		Count(&count)
	return count > 0
}

func (s *ReviewService) SubmitReview(req reviewdto.SubmitReviewRequest) error {
	existing, err := s.reviewRepository.FindReviewByUserAndProduct(s.db, req.UserID, req.ProductID)
	if err != nil {
		return err
	}
	if existing != nil {
		var ce exception.ConflictErrors
		ce.Add("review", "alreadyExist")
		return ce
	}

	review := entity.Review{
		UserID:     req.UserID,
		ProductID:  req.ProductID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		IsVerified: s.isVerifiedBuyer(req.UserID, req.ProductID),
	}
	_, err = s.reviewRepository.CreateReview(s.db, review)
	return err
}

func (s *ReviewService) GetProductReviews(productID uint) ([]reviewdto.ReviewCredential, error) {
	reviews, err := s.reviewRepository.GetProductReviews(s.db, productID)
	if err != nil {
		return nil, err
	}
	var result []reviewdto.ReviewCredential
	for _, r := range reviews {
		cred := reviewdto.ReviewCredential{
			ID:         r.ID,
			UserID:     r.UserID,
			ProductID:  r.ProductID,
			Rating:     r.Rating,
			Comment:    r.Comment,
			IsVerified: r.IsVerified,
			CreatedAt:  r.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
		if r.User.ID != 0 {
			cred.UserPhone = r.User.Phone
		}
		result = append(result, cred)
	}
	return result, nil
}

func (s *ReviewService) DeleteReview(reviewID uint) error {
	review, err := s.reviewRepository.FindReviewByID(s.db, reviewID)
	if err != nil {
		return err
	}
	if review == nil {
		return exception.NotFoundError{Item: "review"}
	}
	return s.reviewRepository.DeleteReviewByID(s.db, reviewID)
}
