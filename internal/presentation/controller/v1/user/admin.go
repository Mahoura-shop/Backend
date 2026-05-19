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

func (userController *AdminUserController) GetUsers(ctx *gin.Context) {
	users, err := userController.userService.GetUsers()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", users)
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

func (userController *AdminUserController) GetUserWallet(ctx *gin.Context) {
	type walletParams struct {
		UserID uint `uri:"userID"`
	}
	params := controller.Validated[walletParams](ctx)

	wallet, err := userController.userService.GetAdminUserWallet(params.UserID)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", wallet)
}

func (userController *AdminUserController) GetSubAdmins(ctx *gin.Context) {
	subAdmins, err := userController.userService.GetSubAdmins()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", subAdmins)
}

func (userController *AdminUserController) CreateSubAdmin(ctx *gin.Context) {
	type params struct {
		Phone  string `json:"phone" validate:"required"`
		RoleID uint   `json:"roleID" validate:"required"`
	}
	p := controller.Validated[params](ctx)
	if err := userController.userService.CreateSubAdmin(p.Phone, p.RoleID); err != nil {
		panic(err)
	}
	controller.Response(ctx, 201, "", nil)
}

func (userController *AdminUserController) AssignSubAdminRole(ctx *gin.Context) {
	type params struct {
		UserID uint `uri:"userID" validate:"required"`
		RoleID uint `json:"roleID" validate:"required"`
	}
	p := controller.Validated[params](ctx)
	if err := userController.userService.AssignSubAdminRole(p.UserID, p.RoleID); err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", nil)
}

func (userController *AdminUserController) RevokeSubAdmin(ctx *gin.Context) {
	type params struct {
		UserID uint `uri:"userID" validate:"required"`
	}
	p := controller.Validated[params](ctx)
	if err := userController.userService.RevokeSubAdmin(p.UserID); err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", nil)
}