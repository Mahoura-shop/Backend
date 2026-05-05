package review

import (
	"strconv"

	"github.com/Mahoura-shop/Backend/bootstrap"
	reviewdto "github.com/Mahoura-shop/Backend/internal/application/dto/review"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerReviewController struct {
	constants     *bootstrap.Constants
	reviewService usecase.ReviewService
}

func NewCustomerReviewController(
	constants *bootstrap.Constants,
	reviewService usecase.ReviewService,
) *CustomerReviewController {
	return &CustomerReviewController{
		constants:     constants,
		reviewService: reviewService,
	}
}

func (c *CustomerReviewController) SubmitReview(ctx *gin.Context) {
	productIDStr := ctx.Param("productID")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		panic(err)
	}

	var req reviewdto.SubmitReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	userID, _ := ctx.Get(c.constants.Context.ID)
	req.UserID = userID.(uint)
	req.ProductID = uint(productID)

	if err := c.reviewService.SubmitReview(req); err != nil {
		panic(err)
	}
	controller.Response(ctx, 201, "", nil)
}

func (c *CustomerReviewController) GetProductReviews(ctx *gin.Context) {
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
