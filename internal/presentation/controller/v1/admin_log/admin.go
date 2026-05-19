package adminlog

import (
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminLogController struct {
	adminLogService usecase.AdminLogService
}

func NewAdminLogController(adminLogService usecase.AdminLogService) *AdminLogController {
	return &AdminLogController{adminLogService: adminLogService}
}

func (c *AdminLogController) GetLogs(ctx *gin.Context) {
	logs, err := c.adminLogService.GetLogs()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", logs)
}
