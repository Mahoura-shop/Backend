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

	products := routerGroup.Group("/products")
	{
		products.GET("", app.Controllers.General.ProductController.GetProducts)
		productSub := products.Group("/:productID")
		{
			productSub.GET("/related", app.Controllers.General.ProductController.GetRelatedProducts)
			productSub.GET("/reviews", app.Controllers.General.ReviewController.GetProductReviews)
		}
		products.GET("/slug/:slug", app.Controllers.General.ProductController.GetProductBySlug)
	}
}