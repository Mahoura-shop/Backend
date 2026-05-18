package httpv1

import (
	"github.com/Mahoura-shop/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	categories := routerGroup.Group("/category")
	{
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
			productsSubGroup.GET("", app.Controllers.Admin.ProductController.GetProduct)
			productsSubGroup.PUT("", app.Controllers.Admin.ProductController.UpdateProduct)
			productsSubGroup.DELETE("", app.Controllers.Admin.ProductController.DeleteProduct)
			productsSubGroup.POST("/images", app.Controllers.Admin.ProductController.AddProductImage)
			imagesSubGroup := productsSubGroup.Group("/images")
			{
				imagesSubGroup.DELETE("/:imageID", app.Controllers.Admin.ProductController.DeleteProductImage)
			}
		}
		categoryProductsSubGroup := products.Group("/category")
		{
			categoryProductsSubGroup.GET("/:categoryID", app.Controllers.Admin.ProductController.GetCategoryProducts)
		}
		products.GET("/prices", app.Controllers.Admin.ProductController.GetProductPrices)
		products.PATCH("/prices", app.Controllers.Admin.ProductController.UpdateProductPrices)
		products.POST("/stock", app.Controllers.Admin.ProductController.UpdateProductStock)
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

	roles := routerGroup.Group("/roles")
	{
		roles.GET("", app.Controllers.Admin.RoleController.GetRoles)
		roles.POST("", app.Controllers.Admin.RoleController.CreateRole)
		roles.GET("/permissions", app.Controllers.Admin.RoleController.GetPermissions)
		rolesSub := roles.Group("/:roleID")
		{
			rolesSub.PUT("", app.Controllers.Admin.RoleController.UpdateRole)
			rolesSub.DELETE("", app.Controllers.Admin.RoleController.DeleteRole)
		}
	}

	admin := routerGroup.Group("/admin")
	{
		admin.GET("/dashboard", app.Controllers.Admin.UserController.GetDashboard)
		admin.GET("/dashboard/orders", app.Controllers.Admin.UserController.GetOrdersChart)
		admin.GET("/dashboard/sales", app.Controllers.Admin.UserController.GetSalesChart)

		upgradeRequests := admin.Group("/upgrade-requests")
		{
			upgradeRequests.GET("", app.Controllers.Admin.UpgradeRequestController.GetUpgradeRequests)
			upgradeRequestsSubGroup := upgradeRequests.Group("/:requestID")
			{
				upgradeRequestsSubGroup.GET("", app.Controllers.Admin.UpgradeRequestController.GetUpgradeRequest)
				upgradeRequestsSubGroup.PATCH("/review", app.Controllers.Admin.UpgradeRequestController.ReviewUpgradeRequest)
			}
		}

		returns := admin.Group("/returns")
		{
			returns.GET("", app.Controllers.Admin.ReturnController.GetReturns)
			returnsSub := returns.Group("/:returnID")
			{
				returnsSub.PATCH("/review", app.Controllers.Admin.ReturnController.ReviewReturn)
				returnsSub.POST("/refund", app.Controllers.Admin.ReturnController.ProcessRefund)
			}
		}

		admin.GET("/contact-messages", app.Controllers.Admin.ContactController.GetAll)
	}

	reviews := routerGroup.Group("/reviews")
	{
		reviewsSub := reviews.Group("/:reviewID")
		{
			reviewsSub.DELETE("", app.Controllers.Admin.ReviewController.DeleteReview)
		}
	}

	users := routerGroup.Group("/users")
	{
		users.GET("", app.Controllers.Admin.UserController.GetUsers)
		usersSubGroup := users.Group("/:userID")
		{
			usersSubGroup.PATCH("/type", app.Controllers.Admin.UpgradeRequestController.ChangeUserType)
			usersSubGroup.GET("/audit-logs", app.Controllers.Admin.UpgradeRequestController.GetUserAuditLogs)
			usersSubGroup.PATCH("/ban", app.Controllers.Admin.UserController.BanUser)
			usersSubGroup.PATCH("/unban", app.Controllers.Admin.UserController.UnbanUser)
			usersSubGroup.GET("/wallet", app.Controllers.Admin.UserController.GetUserWallet)
		}
	}
}
