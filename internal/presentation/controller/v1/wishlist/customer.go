package wishlist

import (
	"strconv"

	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerWishlistController struct {
	constants       *bootstrap.Constants
	wishlistService usecase.WishlistService
}

func NewCustomerWishlistController(
	constants *bootstrap.Constants,
	wishlistService usecase.WishlistService,
) *CustomerWishlistController {
	return &CustomerWishlistController{
		constants:       constants,
		wishlistService: wishlistService,
	}
}

func (c *CustomerWishlistController) GetMyWishlist(ctx *gin.Context) {
	userID, _ := ctx.Get(c.constants.Context.ID)
	items, err := c.wishlistService.GetMyWishlist(userID.(uint))
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", items)
}

func (c *CustomerWishlistController) AddToWishlist(ctx *gin.Context) {
	productIDStr := ctx.Param("productID")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		panic(err)
	}
	userID, _ := ctx.Get(c.constants.Context.ID)
	if err := c.wishlistService.AddToWishlist(userID.(uint), uint(productID)); err != nil {
		panic(err)
	}
	controller.Response(ctx, 201, "", nil)
}

func (c *CustomerWishlistController) RemoveFromWishlist(ctx *gin.Context) {
	productIDStr := ctx.Param("productID")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		panic(err)
	}
	userID, _ := ctx.Get(c.constants.Context.ID)
	if err := c.wishlistService.RemoveFromWishlist(userID.(uint), uint(productID)); err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", nil)
}
