package upgraderequest

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	upgraderequestdto "github.com/Mahoura-shop/Backend/internal/application/dto/upgrade_request"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminUpgradeRequestController struct {
	constants             *bootstrap.Constants
	upgradeRequestService usecase.UpgradeRequestService
}

func NewAdminUpgradeRequestController(
	constants *bootstrap.Constants,
	upgradeRequestService usecase.UpgradeRequestService,
) *AdminUpgradeRequestController {
	return &AdminUpgradeRequestController{
		constants:             constants,
		upgradeRequestService: upgradeRequestService,
	}
}

func (c *AdminUpgradeRequestController) GetUpgradeRequests(ctx *gin.Context) {
	statusStr := ctx.Query("status")
	requests, err := c.upgradeRequestService.GetAllUpgradeRequests(statusStr)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", requests)
}

func (c *AdminUpgradeRequestController) GetUpgradeRequest(ctx *gin.Context) {
	type params struct {
		RequestID uint `uri:"requestID"`
	}
	p := controller.Validated[params](ctx)

	req, err := c.upgradeRequestService.GetUpgradeRequest(p.RequestID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", req)
}

func (c *AdminUpgradeRequestController) ReviewUpgradeRequest(ctx *gin.Context) {
	type params struct {
		RequestID uint `uri:"requestID"`
	}
	p := controller.Validated[params](ctx)

	var req upgraderequestdto.ReviewUpgradeRequestRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	adminID, _ := ctx.Get(c.constants.Context.ID)
	req.AdminID = adminID.(uint)

	if err := c.upgradeRequestService.ReviewUpgradeRequest(p.RequestID, req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, c.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.reviewUpgradeRequest")
	controller.Response(ctx, 200, message, nil)
}

func (c *AdminUpgradeRequestController) ChangeUserType(ctx *gin.Context) {
	type params struct {
		UserID uint `uri:"userID"`
	}
	p := controller.Validated[params](ctx)

	var req upgraderequestdto.ChangeUserTypeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		panic(err)
	}

	adminID, _ := ctx.Get(c.constants.Context.ID)
	req.AdminID = adminID.(uint)

	if err := c.upgradeRequestService.ChangeUserType(p.UserID, req); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, c.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.changeUserType")
	controller.Response(ctx, 200, message, nil)
}

func (c *AdminUpgradeRequestController) GetUserAuditLogs(ctx *gin.Context) {
	type params struct {
		UserID uint `uri:"userID"`
	}
	p := controller.Validated[params](ctx)

	logs, err := c.upgradeRequestService.GetUserAuditLogs(p.UserID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", logs)
}
