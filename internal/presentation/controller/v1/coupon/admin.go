package coupon

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	coupondto "github.com/Mahoura-shop/Backend/internal/application/dto/coupon"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminCouponController struct {
	constants     *bootstrap.Constants
	couponService usecase.CouponService
}

func NewAdminCouponController(
	constants *bootstrap.Constants,
	couponService usecase.CouponService,
) *AdminCouponController {
	return &AdminCouponController{
		constants:     constants,
		couponService: couponService,
	}
}

func (c *AdminCouponController) GetCoupons(ctx *gin.Context) {
	coupons, err := c.couponService.GetCoupons()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", coupons)
}

func (c *AdminCouponController) CreateCoupon(ctx *gin.Context) {
	var req coupondto.CreateCouponRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}
	if err := c.couponService.CreateCoupon(req); err != nil {
		panic(err)
	}
	controller.Response(ctx, 201, "", nil)
}

func (c *AdminCouponController) UpdateCoupon(ctx *gin.Context) {
	type uriParams struct {
		CouponID uint `uri:"couponID" validate:"required"`
	}
	uri := controller.Validated[uriParams](ctx)

	var req coupondto.UpdateCouponRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}
	req.ID = uri.CouponID

	if err := c.couponService.UpdateCoupon(req); err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", nil)
}

func (c *AdminCouponController) DeleteCoupon(ctx *gin.Context) {
	type params struct {
		CouponID uint `uri:"couponID" validate:"required"`
	}
	p := controller.Validated[params](ctx)
	if err := c.couponService.DeleteCoupon(p.CouponID); err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", nil)
}
