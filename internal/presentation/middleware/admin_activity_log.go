package middleware

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/gin-gonic/gin"
)

type AdminActivityLogMiddleware struct {
	constants       *bootstrap.Constants
	adminLogService usecase.AdminLogService
}

func NewAdminActivityLogMiddleware(
	constants *bootstrap.Constants,
	adminLogService usecase.AdminLogService,
) *AdminActivityLogMiddleware {
	return &AdminActivityLogMiddleware{
		constants:       constants,
		adminLogService: adminLogService,
	}
}

func (m *AdminActivityLogMiddleware) LogActivity(ctx *gin.Context) {
	ctx.Next()

	method := ctx.Request.Method
	if method == "GET" || method == "OPTIONS" || method == "HEAD" {
		return
	}

	adminIDRaw, exists := ctx.Get(m.constants.Context.ID)
	if !exists {
		return
	}
	adminID, ok := adminIDRaw.(uint)
	if !ok {
		return
	}

	_ = m.adminLogService.CreateLog(
		adminID,
		method,
		ctx.Request.URL.Path,
		ctx.ClientIP(),
		ctx.Writer.Status(),
	)
}
