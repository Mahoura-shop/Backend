package review

import (
	"strconv"

	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralReviewController struct {
	constants     *bootstrap.Constants
	reviewService usecase.ReviewService
}

func NewGeneralReviewController(
	constants *bootstrap.Constants,
	reviewService usecase.ReviewService,
) *GeneralReviewController {
	return &GeneralReviewController{
		constants:     constants,
		reviewService: reviewService,
	}
}

func (c *GeneralReviewController) GetProductReviews(ctx *gin.Context) {
	productIDStr := ctx.Param("productID")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		panic(err)
	}
	reviews, err := c.reviewService.GetProductReviews(uint(productID))
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", reviews)
}
