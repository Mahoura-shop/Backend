package order

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	orderdto "github.com/Mahoura-shop/Backend/internal/application/dto/order"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerOrderController struct {
	constants    *bootstrap.Constants
	orderService usecase.OrderService
}

func NewCustomerOrderController(
	constants *bootstrap.Constants,
	orderService usecase.OrderService,
) *CustomerOrderController {
	return &CustomerOrderController{
		constants:    constants,
		orderService: orderService,
	}
}

func (orderController *CustomerOrderController) RegisterOrder(ctx *gin.Context) {
	var req orderdto.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	userID, _ := ctx.Get(orderController.constants.Context.ID)

	if err := orderController.orderService.RegisterOrder(userID.(uint), req.PaymentMethod); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, orderController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.registerOrder")
	controller.Response(ctx, 200, message, nil)
}
