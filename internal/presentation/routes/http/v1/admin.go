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
			productsSubGroup.POST("/images", app.Controllers.Admin.ProductController.AddProductImage)
			imagesSubGroup := productsSubGroup.Group("/images")
			{
				imagesSubGroup.DELETE("/:imageID", app.Controllers.Admin.ProductController.DeleteProductImage)
			}
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

	orders := routerGroup.Group("/orders")
	{
		orders.GET("", app.Controllers.Admin.OrderController.GetOrders)
		orderSub := orders.Group("/:orderID")
		{
			orderSub.GET("", app.Controllers.Admin.OrderController.GetOrderDetail)
			orderSub.PATCH("/status", app.Controllers.Admin.OrderController.UpdateOrderStatus)
			orderSub.POST("/cancel", app.Controllers.Admin.OrderController.CancelOrder)
			orderSub.POST("/refund", app.Controllers.Admin.OrderController.FlagRefund)
			orderSub.GET("/instalments", app.Controllers.Admin.OrderController.GetOrderInstalments)
		}
	}

	admin := routerGroup.Group("/admin")
	{
		admin.GET("/dashboard", app.Controllers.Admin.UserController.GetDashboard)
	}

	users := routerGroup.Group("/users")
	{
		users.GET("", app.Controllers.Admin.UserController.GetUsers)
		userSub := users.Group("/:userID")
		{
			userSub.PATCH("/type", app.Controllers.Admin.UpgradeRequestController.ChangeUserType)
			userSub.GET("/audit-logs", app.Controllers.Admin.UpgradeRequestController.GetUserAuditLogs)
		}
	}

	upgradeRequests := routerGroup.Group("/upgrade-requests")
	{
		upgradeRequests.GET("", app.Controllers.Admin.UpgradeRequestController.GetUpgradeRequests)
		requestSub := upgradeRequests.Group("/:requestID")
		{
			requestSub.GET("", app.Controllers.Admin.UpgradeRequestController.GetUpgradeRequest)
			requestSub.PATCH("/review", app.Controllers.Admin.UpgradeRequestController.ReviewUpgradeRequest)
		}
	}
}
