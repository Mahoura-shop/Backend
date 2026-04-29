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

func (c *CustomerOrderController) RegisterOrder(ctx *gin.Context) {
	var req orderdto.CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	userID, _ := ctx.Get(c.constants.Context.ID)
	req.UserID = userID.(uint)

	if err := c.orderService.RegisterOrder(userID.(uint), req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, c.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.registerOrder")
	controller.Response(ctx, 200, message, nil)
}

func (c *CustomerOrderController) GetMyOrders(ctx *gin.Context) {
	userID, _ := ctx.Get(c.constants.Context.ID)
	orders, err := c.orderService.GetUserOrders(userID.(uint))
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", orders)
}

func (c *CustomerOrderController) GetMyOrderDetail(ctx *gin.Context) {
	type params struct {
		OrderID uint `uri:"orderID"`
	}
	p := controller.Validated[params](ctx)

	order, err := c.orderService.GetOrder(p.OrderID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", order)
}

func (c *CustomerOrderController) PayByWallet(ctx *gin.Context) {
	type params struct {
		OrderID uint `uri:"orderID"`
	}
	p := controller.Validated[params](ctx)

	userID, _ := ctx.Get(c.constants.Context.ID)
	if err := c.orderService.PayOrderByWallet(userID.(uint), p.OrderID); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, c.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.payByWallet")
	controller.Response(ctx, 200, message, nil)
}

func (c *CustomerOrderController) InitiatePayment(ctx *gin.Context) {
	type params struct {
		OrderID uint `uri:"orderID"`
	}
	p := controller.Validated[params](ctx)

	userID, _ := ctx.Get(c.constants.Context.ID)
	resp, err := c.orderService.InitiateGatewayPayment(userID.(uint), p.OrderID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", resp)
}

func (c *CustomerOrderController) VerifyPayment(ctx *gin.Context) {
	authority := ctx.Query("Authority")
	status := ctx.Query("Status")

	if err := c.orderService.VerifyGatewayPayment(authority, status); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, c.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.paymentVerified")
	controller.Response(ctx, 200, message, nil)
}

func (c *CustomerOrderController) GetMyOrderInstalments(ctx *gin.Context) {
	type params struct {
		OrderID uint `uri:"orderID"`
	}
	p := controller.Validated[params](ctx)

	instalments, err := c.orderService.GetOrderInstalments(p.OrderID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", instalments)
}
