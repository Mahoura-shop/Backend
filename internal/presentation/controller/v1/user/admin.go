package user

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminUserController struct {
	constants   *bootstrap.Constants
	pagination  *bootstrap.Pagination
	userService usecase.UserService
}

func NewAdminUserController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	userService usecase.UserService,
) *AdminUserController {
	return &AdminUserController{
		constants:   constants,
		pagination:  pagination,
		userService: userService,
	}
}

func (userController *AdminUserController) BanUser(ctx *gin.Context) {
	type banParams struct {
		UserID uint `uri:"userID"`
	}
	params := controller.Validated[banParams](ctx)

	if err := userController.userService.BanUser(params.UserID); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.banUser")
	controller.Response(ctx, 200, message, nil)
}

func (userController *AdminUserController) UnbanUser(ctx *gin.Context) {
	type unbanParams struct {
		UserID uint `uri:"userID"`
	}
	params := controller.Validated[unbanParams](ctx)

	if err := userController.userService.UnbanUser(params.UserID); err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.unbanUser")
	controller.Response(ctx, 200, message, nil)
}