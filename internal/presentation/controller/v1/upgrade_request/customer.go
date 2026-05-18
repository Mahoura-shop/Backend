package upgraderequest

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	upgraderequestdto "github.com/Mahoura-shop/Backend/internal/application/dto/upgrade_request"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerUpgradeRequestController struct {
	constants              *bootstrap.Constants
	upgradeRequestService  usecase.UpgradeRequestService
}

func NewCustomerUpgradeRequestController(
	constants *bootstrap.Constants,
	upgradeRequestService usecase.UpgradeRequestService,
) *CustomerUpgradeRequestController {
	return &CustomerUpgradeRequestController{
		constants:             constants,
		upgradeRequestService: upgradeRequestService,
	}
}

func (c *CustomerUpgradeRequestController) SubmitUpgradeRequest(ctx *gin.Context) {
	req := controller.Validated[upgraderequestdto.SubmitUpgradeRequestRequest](ctx)

	userID, _ := ctx.Get(c.constants.Context.ID)
	req.UserID = userID.(uint)

	if err := c.upgradeRequestService.SubmitUpgradeRequest(req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, c.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.submitUpgradeRequest")
	controller.Response(ctx, 201, message, nil)
}

func (c *CustomerUpgradeRequestController) GetMyUpgradeRequests(ctx *gin.Context) {
	userID, _ := ctx.Get(c.constants.Context.ID)

	requests, err := c.upgradeRequestService.GetMyUpgradeRequests(userID.(uint))
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", requests)
}
