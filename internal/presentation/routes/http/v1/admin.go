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

	currencies := routerGroup.Group("/currency")
	{
		currencies.GET("", app.Controllers.Admin.CurrencyController.GetCurrencies)
		currenciesSubGroup := currencies.Group("/:currencyID") 
		{
			currenciesSubGroup.PUT("", app.Controllers.Admin.CurrencyController.UpdateCurrency)
		}
	}
	
	brands := routerGroup.Group("/brand")
	{
		brands.GET("", app.Controllers.Admin.BrandController.GetBrands)
		brands.POST("", app.Controllers.Admin.BrandController.CreateBrand)
		brandsSubGroup := brands.Group("/:brandID") 
		{
			brandsSubGroup.PUT("", app.Controllers.Admin.BrandController.UpdateBrand)
			brandsSubGroup.DELETE("", app.Controllers.Admin.BrandController.DeleteBrand)
		}
	}

	products := routerGroup.Group("/product")
	{
		products.POST("", app.Controllers.Admin.ProductController.CreateProduct)
		products.GET("", app.Controllers.Admin.ProductController.GetProducts)
		productsSubGroup := products.Group("/:productID") 
		{
			productsSubGroup.PUT("", app.Controllers.Admin.ProductController.UpdateProduct)
			productsSubGroup.DELETE("", app.Controllers.Admin.ProductController.DeleteProduct)
		}
		productsSlugSubgroup := products.Group("/:slug") 
		{
			productsSlugSubgroup.GET("", app.Controllers.Admin.ProductController.GetProduct)
		}
		categoryProductsSubGroup := products.Group("/category")
		{
			categoryProductsSubGroup.GET("/:categoryID", app.Controllers.Admin.ProductController.GetCategoryProducts)
		}
		products.GET("/prices", app.Controllers.Admin.ProductController.GetProductPrices)
		products.PATCH("/prices", app.Controllers.Admin.ProductController.UpdateProductPrices)
	}

	admin := routerGroup.Group("/admin")
	{
		admin.POST("/login", app.Controllers.General.UserController.AdminLogin)
		admin.GET("/dashboard", app.Controllers.Admin.UserController.GetDashboard)
	}
}
