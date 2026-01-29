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
	
	brands := routerGroup.Group("/brand")
	{
		// brands.GET("", app.Controllers.Admin.BrandController.GetBrands)
		brands.POST("", app.Controllers.Admin.BrandController.CreateBrand)
		// brandsSubGroup := brands.Group("/:brandID") 
		// {
		// 	brandsSubGroup.PUT("", app.Controllers.Admin.BrandController.UpdateBrand)
		// 	brandsSubGroup.DELETE("", app.Controllers.Admin.BrandController.DeleteBrand)
		// }
	}

	products := routerGroup.Group("/product")
	{
		products.POST("", app.Controllers.Admin.ProductController.CreateProduct)
		products.GET("", app.Controllers.Admin.ProductController.GetProducts)
		productsSubGroup := products.Group("/:productID") 
		{
			productsSubGroup.GET("", app.Controllers.Admin.ProductController.GetProduct)
			productsSubGroup.PUT("", app.Controllers.Admin.ProductController.UpdateProduct)
			productsSubGroup.DELETE("", app.Controllers.Admin.ProductController.DeleteProduct)
		}
	}

	admin := routerGroup.Group("/admin")
	{
		admin.POST("/login", app.Controllers.General.UserController.AdminLogin)
	}
}
