package review

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminReviewController struct {
	constants     *bootstrap.Constants
	reviewService usecase.ReviewService
}

func NewAdminReviewController(
	constants *bootstrap.Constants,
	reviewService usecase.ReviewService,
) *AdminReviewController {
	return &AdminReviewController{
		constants:     constants,
		reviewService: reviewService,
	}
}

func (c *AdminReviewController) DeleteReview(ctx *gin.Context) {
	type params struct {
		ReviewID uint `uri:"reviewID" validate:"required"`
	}
	p := controller.Validated[params](ctx)

	if err := c.reviewService.DeleteReview(p.ReviewID); err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", nil)
}
