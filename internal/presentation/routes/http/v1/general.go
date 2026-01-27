package httpv1

import (
	"github.com/Mahoura-shop/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupGeneralRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	test := routerGroup.Group("/test")
	{
		test.POST("/", app.Controllers.General.TestController.Test)
		// test.POST("/", func(c *gin.Context) {
		// 	c.JSON(200, gin.H{"ok": true})
		// })
	}
}
