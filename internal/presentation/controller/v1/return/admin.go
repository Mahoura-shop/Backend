package returnctrl

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	returndto "github.com/Mahoura-shop/Backend/internal/application/dto/return"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminReturnController struct {
	constants     *bootstrap.Constants
	returnService usecase.ReturnService
}

func NewAdminReturnController(
	constants *bootstrap.Constants,
	returnService usecase.ReturnService,
) *AdminReturnController {
	return &AdminReturnController{
		constants:     constants,
		returnService: returnService,
	}
}

func (c *AdminReturnController) GetReturns(ctx *gin.Context) {
	status := ctx.Query("status")
	returns, err := c.returnService.GetAllReturns(status)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", returns)
}

func (c *AdminReturnController) ReviewReturn(ctx *gin.Context) {
	type params struct {
		ReturnID uint `uri:"returnID"`
	}
	p := controller.Validated[params](ctx)

	var req returndto.ReviewReturnRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}
	adminID, _ := ctx.Get(c.constants.Context.ID)
	req.AdminID = adminID.(uint)

	if err := c.returnService.ReviewReturn(p.ReturnID, req); err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "return reviewed", nil)
}

func (c *AdminReturnController) ProcessRefund(ctx *gin.Context) {
	type params struct {
		ReturnID uint `uri:"returnID"`
	}
	p := controller.Validated[params](ctx)

	adminID, _ := ctx.Get(c.constants.Context.ID)
	if err := c.returnService.ProcessRefund(p.ReturnID, adminID.(uint)); err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "refund processed", nil)
}
