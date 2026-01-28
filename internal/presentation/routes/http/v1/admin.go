package httpv1

import (
	"github.com/Mahoura-shop/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	categories := routerGroup.Group("/category")
	{
		categories.GET("", app.Controllers.Admin.CategoryController.GetCategories)
		categories.POST("", app.Controllers.Admin.CategoryController.CreateCategory)
		categoriesSubGroup := categories.Group("/:categoryID") 
		{
			categoriesSubGroup.PUT("", app.Controllers.Admin.CategoryController.UpdateCategory)
			categoriesSubGroup.DELETE("", app.Controllers.Admin.CategoryController.DeleteCategory)
		}
	}

	products := routerGroup.Group("/product")
	{
		products.POST("", app.Controllers.Admin.ProductController.CreateProduct)
		// products.GET("", app.Controllers.Admin.CategoryController.GetCategorie)
		// products.POST("/buy", app.Controllers.Admin.CategoryController.GetCategGetCategorieories)
		// products.POST("/sell", app.Controllers.Admin.CategoryController.GetCategorie)
		// productsSubGroup := categories.Group("/:productID") 
		// {
		// 	productsSubGroup.DELETE("", app.Controllers.Admin.CategoryController.GetCategorie)
		// 	productsSubGroup.PUT("", app.Controllers.Admin.CategoryController.GetCategorie)
		// 	productsSubGroup.GET("", app.Controllers.Admin.CategoryController.GetCategorie)
		// }
	}

	admin := routerGroup.Group("/admin")
	{
		admin.POST("/login", app.Controllers.General.UserController.AdminLogin)
	}
}
