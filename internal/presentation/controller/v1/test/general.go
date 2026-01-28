package test

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	testdto "github.com/Mahoura-shop/Backend/internal/application/dto/test"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralTestController struct {
	constants      *bootstrap.Constants
	testService     usecase.TestService
}

func NewGeneralTestController(
	constants *bootstrap.Constants,
	testService usecase.TestService,
) *GeneralTestController {
	return &GeneralTestController{
		constants:   constants,
		testService: testService,
	}
}

func (testController *GeneralTestController) Test(ctx *gin.Context) {
	type registerParams struct {
		Test string `json:"test" validate:"required"`
	}
	params := controller.Validated[registerParams](ctx)
	testInfo := testdto.BasicTestRequest{
		Test: params.Test,
	}
	res, err := testController.testService.Test(testInfo.Test); 
	if (err != nil) {
		panic(err)
	}

	trans := controller.GetTranslator(ctx, testController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.test")
	controller.Response(ctx, 200, message, res)
}