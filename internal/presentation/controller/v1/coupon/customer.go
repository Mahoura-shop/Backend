package coupon

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	coupondto "github.com/Mahoura-shop/Backend/internal/application/dto/coupon"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerCouponController struct {
	constants     *bootstrap.Constants
	couponService usecase.CouponService
}

func NewCustomerCouponController(
	constants *bootstrap.Constants,
	couponService usecase.CouponService,
) *CustomerCouponController {
	return &CustomerCouponController{
		constants:     constants,
		couponService: couponService,
	}
}

func (c *CustomerCouponController) ValidateCoupon(ctx *gin.Context) {
	var req coupondto.ValidateCouponRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}
	result, err := c.couponService.ValidateCoupon(req.Code, req.IRRPrice)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", result)
}
