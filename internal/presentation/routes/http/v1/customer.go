package httpv1

import (
	"github.com/Mahoura-shop/Backend/wire"
	"github.com/gin-gonic/gin"
)

func SetupCustomerRoutes(routerGroup *gin.RouterGroup, app *wire.Application) {
	profile := routerGroup.Group("/profile")
	{
		profile.GET("", app.Controllers.Customer.UserController.GetMyProfile)
	}

	wallet := routerGroup.Group("/wallet")
	{
		wallet.GET("", app.Controllers.Customer.UserController.GetUserWalletBalance)
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
}
