package order

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	orderdto "github.com/Mahoura-shop/Backend/internal/application/dto/order"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminOrderController struct {
	constants    *bootstrap.Constants
	orderService usecase.OrderService
}

func NewAdminOrderController(
	constants *bootstrap.Constants,
	orderService usecase.OrderService,
) *AdminOrderController {
	return &AdminOrderController{
		constants:    constants,
		orderService: orderService,
	}
}

func (c *AdminOrderController) GetOrders(ctx *gin.Context) {
	statusStr := ctx.Query("status")
	if statusStr != "" {
		// Filter by status if provided
		var status enum.OrderStatus
		switch statusStr {
		case "pending":
			status = enum.OrderStatusPending
		case "paid":
			status = enum.OrderStatusPaid
		case "shipped":
			status = enum.OrderStatusShipped
		case "cancelled":
			status = enum.OrderStatusCancelled
		}
		if status > 0 {
			orders, err := c.orderService.GetOrdersByStatus(status)
			if err != nil {
				panic(err)
			}
			controller.Response(ctx, 200, "", orders)
			return
		}
	}

	orders, err := c.orderService.GetOrders()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", orders)
}

func (c *AdminOrderController) GetOrderDetail(ctx *gin.Context) {
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

func (c *AdminOrderController) UpdateOrderStatus(ctx *gin.Context) {
	type params struct {
		OrderID uint `uri:"orderID"`
	}
	p := controller.Validated[params](ctx)

	var req orderdto.UpdateOrderStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	adminID, _ := ctx.Get(c.constants.Context.ID)
	req.ChangedByID = adminID.(uint)

	if err := c.orderService.UpdateOrderStatus(p.OrderID, req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, c.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateOrderStatus")
	controller.Response(ctx, 200, message, nil)
}

func (c *AdminOrderController) CancelOrder(ctx *gin.Context) {
	type params struct {
		OrderID uint `uri:"orderID"`
	}
	type body struct {
		Reason string `json:"reason" binding:"required"`
	}
	p := controller.Validated[params](ctx)

	var b body
	if err := ctx.ShouldBindJSON(&b); err != nil {
		panic(err)
	}

	if err := c.orderService.CancelOrder(p.OrderID, b.Reason); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, c.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.cancelOrder")
	controller.Response(ctx, 200, message, nil)
}

func (c *AdminOrderController) FlagRefund(ctx *gin.Context) {
	type params struct {
		OrderID uint `uri:"orderID"`
	}
	p := controller.Validated[params](ctx)

	if err := c.orderService.FlagOrderRefund(p.OrderID); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, c.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.flagRefund")
	controller.Response(ctx, 200, message, nil)
}

func (c *AdminOrderController) GetOrderInstalments(ctx *gin.Context) {
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
