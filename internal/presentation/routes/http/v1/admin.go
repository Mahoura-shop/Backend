package httpv1

import (
	"github.com/Mahoura-shop/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	category := routerGroup.Group("/category")
	{
		category.GET("", app.Controllers.Admin.CategoryController.GetCategories)
		category.POST("", app.Controllers.Admin.CategoryController.CreateCategory)
		// category.PUT("/", app.Controllers.Admin.TestController.Test)
		// category.DELETE("/", app.Controllers.Admin.TestController.Test)
	}

	admin := routerGroup.Group("/admin")
	{
		admin.POST("/login", app.Controllers.General.UserController.AdminLogin)
	}
}
