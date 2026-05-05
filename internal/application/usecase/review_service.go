package usecase

import reviewdto "github.com/Mahoura-shop/Backend/internal/application/dto/review"

type ReviewService interface {
	SubmitReview(reviewdto.SubmitReviewRequest) error
	GetProductReviews(productID uint) ([]reviewdto.ReviewCredential, error)
	DeleteReview(reviewID uint) error
}
