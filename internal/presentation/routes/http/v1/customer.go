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
}
