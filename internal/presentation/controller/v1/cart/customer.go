package cart

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	cartdto "github.com/Mahoura-shop/Backend/internal/application/dto/cart"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerCartController struct {
	constants   *bootstrap.Constants
	cartService usecase.CartService
}

func NewCustomerCartController(
	constants *bootstrap.Constants,
	cartService usecase.CartService,
) *CustomerCartController {
	return &CustomerCartController{
		constants:   constants,
		cartService: cartService,
	}
}

func (cartController *CustomerCartController) GetUserCart(ctx *gin.Context) {
	userID, _ := ctx.Get(cartController.constants.Context.ID)
	cart, err := cartController.cartService.GetUserCart(userID.(uint))
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", cart)
}

func (cartController *CustomerCartController) AddProductToCart(ctx *gin.Context) {
	userID, _ := ctx.Get(cartController.constants.Context.ID)
	type addProductToCartParams struct {
		ProductID uint `uri:"productID" validate:"required"`
	}
	params := controller.Validated[addProductToCartParams](ctx)
	addProductToCartInfo := cartdto.UpdateProductCountInCart{
		UserID: userID.(uint),
		ProductID: params.ProductID,
	}

	err := cartController.cartService.AddProductToCart(addProductToCartInfo); 
	if err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, cartController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.addProductToCart")
	controller.Response(ctx, 200, message, nil)
}

func (cartController *CustomerCartController) RemoveProductFromCart(ctx *gin.Context) {
	userID, _ := ctx.Get(cartController.constants.Context.ID)
	type removeProductFromCartParams struct {
		ProductID uint `uri:"productID" validate:"required"`
	}
	params := controller.Validated[removeProductFromCartParams](ctx)
	addProductToCartInfo := cartdto.UpdateProductCountInCart{
		UserID: userID.(uint),
		ProductID: params.ProductID,
	}

	err := cartController.cartService.RemoveProductFromCart(addProductToCartInfo); 
	if err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, cartController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.removeProductFromCart")
	controller.Response(ctx, 200, message, nil)
}