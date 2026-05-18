package returnctrl

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	returndto "github.com/Mahoura-shop/Backend/internal/application/dto/return"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerReturnController struct {
	constants     *bootstrap.Constants
	returnService usecase.ReturnService
}

func NewCustomerReturnController(
	constants *bootstrap.Constants,
	returnService usecase.ReturnService,
) *CustomerReturnController {
	return &CustomerReturnController{
		constants:     constants,
		returnService: returnService,
	}
}

func (c *CustomerReturnController) RequestReturn(ctx *gin.Context) {
	var req returndto.RequestReturnRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}
	userID, _ := ctx.Get(c.constants.Context.ID)
	req.UserID = userID.(uint)

	if err := c.returnService.RequestReturn(req); err != nil {
		panic(err)
	}
	controller.Response(ctx, 201, "return request submitted", nil)
}

func (c *CustomerReturnController) GetMyReturns(ctx *gin.Context) {
	userID, _ := ctx.Get(c.constants.Context.ID)
	returns, err := c.returnService.GetMyReturns(userID.(uint))
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", returns)
}
