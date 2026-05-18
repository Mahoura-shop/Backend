package httpv1

import (
	"github.com/Mahoura-shop/Backend/internal/application/service/shipping"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
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
	}

	province := routerGroup.Group("/province")
	{
		province.GET("", app.Controllers.General.AddressController.GetProvince)
		province.GET("/:provinceID/cities", app.Controllers.General.AddressController.GetProvinceCities)
	}

	category := routerGroup.Group("/category")
	{
		category.GET("", app.Controllers.Admin.CategoryController.GetCategories)
	}

	brand := routerGroup.Group("/brand")
	{
		brand.GET("", app.Controllers.Admin.BrandController.GetBrands)
	}

	routerGroup.GET("/shipping", func(ctx *gin.Context) {
		controller.Response(ctx, 200, "", gin.H{"shippingCost": shipping.DefaultShippingCost})
	})

	contact := routerGroup.Group("/contact")
	{
		contact.POST("", app.Controllers.General.ContactController.Submit)
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