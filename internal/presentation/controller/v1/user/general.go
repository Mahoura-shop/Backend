package user

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralUserController struct {
	constants   *bootstrap.Constants
	userService usecase.UserService
	jwtService  usecase.JWTService
}

func NewGeneralUserController(
	constants *bootstrap.Constants,
	userService usecase.UserService,
	jwtService usecase.JWTService,
) *GeneralUserController {
	return &GeneralUserController{
		constants:   constants,
		userService: userService,
		jwtService:  jwtService,
	}
}

func (userController *GeneralUserController) Auth(ctx *gin.Context) {
	type authParams struct {
		Phone string `json:"phone" validate:"required"`
	}
	params := controller.Validated[authParams](ctx)
	registerInfo := userdto.AuthRequest{
		Phone: params.Phone,
	}
	if err := userController.userService.Auth(registerInfo); err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.auth")
	controller.Response(ctx, 200, message, nil)
}

func (userController *GeneralUserController) VerifyAuth(ctx *gin.Context) {
	type verifyPhoneParams struct {
		Phone string `json:"phone" validate:"required"`
		OTP   string `json:"otp" validate:"required"`
	}
	params := controller.Validated[verifyPhoneParams](ctx)
	verifyOTPInfo := userdto.VerifyAuthRequest{
		Phone: params.Phone,
		OTP:   params.OTP,
	}
	user, err := userController.userService.VerifyAuth(verifyOTPInfo); 
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.verifyAuth")
	controller.Response(ctx, 200, message, user)
}

func (userController *GeneralUserController) RefreshToken(ctx *gin.Context) {
	type refreshTokenParams struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}
	params := controller.Validated[refreshTokenParams](ctx)
	claims, err := userController.jwtService.ValidateToken(params.RefreshToken)
	if err != nil {
		panic(err)
	}

	userID := uint(claims["sub"].(float64))
	accessToken, _, err := userController.jwtService.GenerateToken(userID)
	if err != nil {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, userController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.refreshToken")
	controller.Response(ctx, 200, message, accessToken)
}


func (userController *AdminUserController) GetDashboard(ctx *gin.Context) {
	dashboard, err := userController.userService.GetDashboard()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", dashboard)
}

func (userController *AdminUserController) GetOrdersChart(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "week")
	data, err := userController.userService.GetOrdersChart(period)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", data)
}

func (userController *AdminUserController) GetSalesChart(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "week")
	data, err := userController.userService.GetSalesChart(period)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", data)
}

func (userController *AdminUserController) GetVisitsChart(ctx *gin.Context) {
	period := ctx.DefaultQuery("period", "week")
	data, err := userController.userService.GetVisitsChart(period)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", data)
}

func (userController *AdminUserController) GetProductVisitsChart(ctx *gin.Context) {
	type params struct {
		ProductID uint `uri:"productID" validate:"required"`
	}
	p := controller.Validated[params](ctx)
	period := ctx.DefaultQuery("period", "week")
	data, err := userController.userService.GetProductVisitsChart(p.ProductID, period)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", data)
}

func (userController *AdminUserController) GetProductOrdersChart(ctx *gin.Context) {
	type params struct {
		ProductID uint `uri:"productID" validate:"required"`
	}
	p := controller.Validated[params](ctx)
	period := ctx.DefaultQuery("period", "week")
	data, err := userController.userService.GetProductOrdersChart(p.ProductID, period)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", data)
}