package user

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
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