package user

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerUserController struct {
	constants   *bootstrap.Constants
	userService usecase.UserService
}

func NewCustomerUserController(
	constants *bootstrap.Constants,
	userService usecase.UserService,
) *CustomerUserController {
	return &CustomerUserController{
		constants:   constants,
		userService: userService,
	}
}

func (userController *CustomerUserController) GetMyProfile(ctx *gin.Context) {
	userID, _ := ctx.Get(userController.constants.Context.ID)
	profile, err := userController.userService.GetUserCredential(userID.(uint))
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", profile)
}

func (userController *CustomerUserController) GetUserWalletBalance(ctx *gin.Context) {
	userID, _ := ctx.Get(userController.constants.Context.ID)
	balance, err := userController.userService.GetUserWalletBalance(userID.(uint))
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", balance)
}

func (userController *CustomerUserController) DepositWallet(ctx *gin.Context) {
	userID, _ := ctx.Get(userController.constants.Context.ID)
	type DepositWalletParams struct {
		Amount uint `json:"amount" validate:"required"`
	}
	params := controller.Validated[DepositWalletParams](ctx)
	
	balanceUpdateInfo := userdto.UserBalanceUpdate{
		UserID: userID.(uint),
		Amount: params.Amount,
	}

	newBalance, err := userController.userService.DepositWallet(balanceUpdateInfo); 
	if err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deposit")
	controller.Response(ctx, 200, message, newBalance)
}

func (userController *CustomerUserController) WithdrawWallet(ctx *gin.Context) {
	userID, _ := ctx.Get(userController.constants.Context.ID)
	type WithdrawWalletParams struct {
		Amount uint `json:"amount" validate:"required"`
	}
	params := controller.Validated[WithdrawWalletParams](ctx)
	
	balanceUpdateInfo := userdto.UserBalanceUpdate{
		UserID: userID.(uint),
		Amount: params.Amount,
	}

	newBalance, err := userController.userService.WithdrawWallet(balanceUpdateInfo); 
	if err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.withdraw")
	controller.Response(ctx, 200, message, newBalance)
}