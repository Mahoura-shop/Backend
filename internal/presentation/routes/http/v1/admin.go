package httpv1

import (
	"github.com/Mahoura-shop/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupAdminRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	perm := app.Middlewares.Authentication.RequirePermission

	categories := routerGroup.Group("/category")
	{
		categories.POST("", perm("category:create"), app.Controllers.Admin.CategoryController.CreateCategory)
		categoriesSubGroup := categories.Group("/:categoryID")
		{
			categoriesSubGroup.PUT("", perm("category:edit"), app.Controllers.Admin.CategoryController.UpdateCategory)
			categoriesSubGroup.DELETE("", perm("category:delete"), app.Controllers.Admin.CategoryController.DeleteCategory)
		}
	}

	currencies := routerGroup.Group("/currency")
	{
		currencies.GET("", app.Controllers.Admin.CurrencyController.GetCurrencies)
		currenciesSubGroup := currencies.Group("/:currencyID")
		{
			currenciesSubGroup.PUT("", perm("update:currencies"), app.Controllers.Admin.CurrencyController.UpdateCurrency)
		}
	}

	brands := routerGroup.Group("/brand")
	{
		brands.POST("", perm("brand:create"), app.Controllers.Admin.BrandController.CreateBrand)
		brandsSubGroup := brands.Group("/:brandID")
		{
			brandsSubGroup.PUT("", perm("brand:edit"), app.Controllers.Admin.BrandController.UpdateBrand)
			brandsSubGroup.DELETE("", perm("brand:delete"), app.Controllers.Admin.BrandController.DeleteBrand)
		}
	}

	products := routerGroup.Group("/product")
	{
		products.POST("", perm("product:create"), app.Controllers.Admin.ProductController.CreateProduct)
		products.GET("", perm("product:see"), app.Controllers.Admin.ProductController.GetProducts)
		productsSubGroup := products.Group("/:productID")
		{
			productsSubGroup.GET("", perm("product:see"), app.Controllers.Admin.ProductController.GetProduct)
			productsSubGroup.PUT("", perm("product:edit"), app.Controllers.Admin.ProductController.UpdateProduct)
			productsSubGroup.DELETE("", perm("product:delete"), app.Controllers.Admin.ProductController.DeleteProduct)
			productsSubGroup.POST("/images", perm("product:edit"), app.Controllers.Admin.ProductController.AddProductImage)
			imagesSubGroup := productsSubGroup.Group("/images")
			{
				imagesSubGroup.DELETE("/:imageID", perm("product:edit"), app.Controllers.Admin.ProductController.DeleteProductImage)
			}
		}
		categoryProductsSubGroup := products.Group("/category")
		{
			categoryProductsSubGroup.GET("/:categoryID", perm("product:see"), app.Controllers.Admin.ProductController.GetCategoryProducts)
		}
		brandProductsSubGroup := products.Group("/brand")
		{
			brandProductsSubGroup.GET("/:brandID", perm("product:see"), app.Controllers.Admin.ProductController.GetBrandProducts)
		}
		products.GET("/prices", perm("product:batch_price"), app.Controllers.Admin.ProductController.GetProductPrices)
		products.PATCH("/prices", perm("product:batch_price"), app.Controllers.Admin.ProductController.UpdateProductPrices)
		products.POST("/stock", perm("product:batch_inventory"), app.Controllers.Admin.ProductController.UpdateProductStock)
	}

	orders := routerGroup.Group("/orders")
	{
		orders.GET("", perm("order:see"), app.Controllers.Admin.OrderController.GetOrders)
		orderSub := orders.Group("/:orderID")
		{
			orderSub.GET("", perm("order:see"), app.Controllers.Admin.OrderController.GetOrderDetail)
			orderSub.PATCH("/status", perm("order:update_status"), app.Controllers.Admin.OrderController.UpdateOrderStatus)
			orderSub.POST("/cancel", perm("order:cancel"), app.Controllers.Admin.OrderController.CancelOrder)
			orderSub.POST("/refund", perm("order:cancel"), app.Controllers.Admin.OrderController.FlagRefund)
			orderSub.GET("/instalments", perm("order:see"), app.Controllers.Admin.OrderController.GetOrderInstalments)
		}
	}

	roles := routerGroup.Group("/roles")
	{
		roles.GET("", perm("rbac:see"), app.Controllers.Admin.RoleController.GetRoles)
		roles.POST("", perm("rbac:create"), app.Controllers.Admin.RoleController.CreateRole)
		roles.GET("/permissions", perm("rbac:see"), app.Controllers.Admin.RoleController.GetPermissions)
		rolesSub := roles.Group("/:roleID")
		{
			rolesSub.PUT("", perm("rbac:edit"), app.Controllers.Admin.RoleController.UpdateRole)
			rolesSub.DELETE("", perm("rbac:delete"), app.Controllers.Admin.RoleController.DeleteRole)
		}
	}

	admin := routerGroup.Group("/admin")
	{
		admin.GET("/logs", perm("adminlogs:see"), app.Controllers.Admin.AdminLogController.GetLogs)
		admin.GET("/dashboard", app.Controllers.Admin.UserController.GetDashboard)
		admin.GET("/dashboard/orders", app.Controllers.Admin.UserController.GetOrdersChart)
		admin.GET("/dashboard/sales", app.Controllers.Admin.UserController.GetSalesChart)
		admin.GET("/dashboard/visits", app.Controllers.Admin.UserController.GetVisitsChart)

		adminProducts := admin.Group("/products")
		{
			adminProductsSub := adminProducts.Group("/:productID")
			{
				adminProductsSub.GET("/visits", perm("product:see"), app.Controllers.Admin.UserController.GetProductVisitsChart)
				adminProductsSub.GET("/orders", perm("product:see"), app.Controllers.Admin.UserController.GetProductOrdersChart)
			}
		}

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

		admin.GET("/contact-messages", perm("contact:see"), app.Controllers.Admin.ContactController.GetAll)
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
		users.GET("", perm("users:see"), app.Controllers.Admin.UserController.GetUsers)
		usersSubGroup := users.Group("/:userID")
		{
			usersSubGroup.PATCH("/type", perm("users:edit_role"), app.Controllers.Admin.UpgradeRequestController.ChangeUserType)
			usersSubGroup.GET("/audit-logs", app.Controllers.Admin.UpgradeRequestController.GetUserAuditLogs)
			usersSubGroup.PATCH("/ban", perm("users:ban"), app.Controllers.Admin.UserController.BanUser)
			usersSubGroup.PATCH("/unban", perm("users:unban"), app.Controllers.Admin.UserController.UnbanUser)
			usersSubGroup.GET("/wallet", perm("users:wallet"), app.Controllers.Admin.UserController.GetUserWallet)
		}
	}

	subAdmins := routerGroup.Group("/subadmins")
	{
		subAdmins.GET("", app.Controllers.Admin.UserController.GetSubAdmins)
		subAdmins.POST("", app.Controllers.Admin.UserController.CreateSubAdmin)
		subAdminsSub := subAdmins.Group("/:userID")
		{
			subAdminsSub.PATCH("/role", app.Controllers.Admin.UserController.AssignSubAdminRole)
			subAdminsSub.DELETE("", app.Controllers.Admin.UserController.RevokeSubAdmin)
		}
	}
}
