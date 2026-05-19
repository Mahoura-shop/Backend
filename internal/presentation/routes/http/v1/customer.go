package httpv1

import (
	"github.com/Mahoura-shop/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	profile := routerGroup.Group("/profile")
	{
		profile.GET("", app.Controllers.Customer.UserController.GetMyProfile)
		profile.PATCH("", app.Controllers.Customer.UserController.UpdateMyProfile)
	}

	wallet := routerGroup.Group("/wallet")
	{
		wallet.GET("", app.Controllers.Customer.UserController.GetUserWalletBalance)
		wallet.GET("/history", app.Controllers.Customer.UserController.GetWalletHistory)
		wallet.POST("/deposit", app.Controllers.Customer.UserController.DepositWallet)
		wallet.POST("/withdraw", app.Controllers.Customer.UserController.WithdrawWallet)
	}

	cart := routerGroup.Group("/cart")
	{
		cart.GET("", app.Controllers.Customer.CartController.GetUserCart)
		productActionsSubgroup := cart.Group("/:productID")
		{
			productActionsSubgroup.POST("/add", app.Controllers.Customer.CartController.AddProductToCart)
			productActionsSubgroup.POST("/remove", app.Controllers.Customer.CartController.RemoveProductFromCart)
		}
	}

	order := routerGroup.Group("/order")
	{
		order.POST("", app.Controllers.Customer.OrderController.RegisterOrder)
		order.GET("", app.Controllers.Customer.OrderController.GetMyOrders)
		order.GET("/pay/verify", app.Controllers.Customer.OrderController.VerifyPayment)
		orderSub := order.Group("/:orderID")
		{
			orderSub.GET("", app.Controllers.Customer.OrderController.GetMyOrderDetail)
			orderSub.POST("/pay/wallet", app.Controllers.Customer.OrderController.PayByWallet)
			orderSub.POST("/pay/gateway", app.Controllers.Customer.OrderController.InitiatePayment)
			orderSub.GET("/instalments", app.Controllers.Customer.OrderController.GetMyOrderInstalments)
		}
	}

	address := routerGroup.Group("/address")
	{
		address.GET("", app.Controllers.Customer.AddressController.GetCustomerAddresses)
		address.POST("", app.Controllers.Customer.AddressController.CreateUserAddress)
	}

	upgradeRequests := routerGroup.Group("/upgrade-requests")
	{
		upgradeRequests.POST("", app.Controllers.Customer.UpgradeRequestController.SubmitUpgradeRequest)
		upgradeRequests.GET("", app.Controllers.Customer.UpgradeRequestController.GetMyUpgradeRequests)
	}

	wishlist := routerGroup.Group("/wishlist")
	{
		wishlist.GET("", app.Controllers.Customer.WishlistController.GetMyWishlist)
		wishlistSub := wishlist.Group("/:productID")
		{
			wishlistSub.POST("", app.Controllers.Customer.WishlistController.AddToWishlist)
			wishlistSub.DELETE("", app.Controllers.Customer.WishlistController.RemoveFromWishlist)
		}
	}

	products := routerGroup.Group("/products")
	{
		productSub := products.Group("/:productID")
		{
			productSub.POST("/review", app.Controllers.Customer.ReviewController.SubmitReview)
		}
	}

	reviews := routerGroup.Group("/reviews")
	{
		reviews.GET("", app.Controllers.Customer.ReviewController.GetMyReviews)
	}

	coupon := routerGroup.Group("/coupon")
	{
		coupon.POST("/validate", app.Controllers.Customer.CouponController.ValidateCoupon)
	}

	returns := routerGroup.Group("/returns")
	{
		returns.POST("", app.Controllers.Customer.ReturnController.RequestReturn)
		returns.GET("", app.Controllers.Customer.ReturnController.GetMyReturns)
	}

	notifications := routerGroup.Group("/notifications")
	{
		notifications.GET("", app.Controllers.Customer.NotificationController.GetNotifications)
		notifications.GET("/unread-count", app.Controllers.Customer.NotificationController.CountUnread)
		notifications.PATCH("/read-all", app.Controllers.Customer.NotificationController.MarkAllAsRead)
		notifications.PATCH("/:id/read", app.Controllers.Customer.NotificationController.MarkAsRead)
	}
}
