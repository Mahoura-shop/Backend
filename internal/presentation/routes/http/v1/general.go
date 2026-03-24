package httpv1

import (
	"github.com/Mahoura-shop/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	test := routerGroup.Group("/test")
	{
		test.POST("", app.Controllers.General.TestController.Test)
	}

	auth := routerGroup.Group("/auth")
	{
		auth.POST("", app.Controllers.General.UserController.Auth)
		auth.POST("/verify", app.Controllers.General.UserController.VerifyAuth)
		auth.POST("/login", app.Controllers.General.UserController.AdminLogin)
	}
}