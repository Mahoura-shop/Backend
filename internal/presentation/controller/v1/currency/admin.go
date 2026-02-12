package currency

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	currencydto "github.com/Mahoura-shop/Backend/internal/application/dto/currency"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminCurrencyController struct {
	constants       *bootstrap.Constants
	pagination      *bootstrap.Pagination
	currencyService usecase.CurrencyService
}

func NewAdminCurrencyController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	currencyService usecase.CurrencyService,
) *AdminCurrencyController {
	return &AdminCurrencyController{
		constants:       constants,
		pagination:      pagination,
		currencyService: currencyService,
	}
}

func (currencyController *AdminCurrencyController) GetCurrencies(ctx *gin.Context) {
	currencies, err := currencyController.currencyService.GetCurrencies(); 
	if (err != nil) {
		panic(err)
	}
	controller.Response(ctx, 200, "", currencies)
}

func (currencyController *AdminCurrencyController) UpdateCurrency(ctx *gin.Context) {
	type updateCurrencyParams struct {
		ID          uint  `uri:"currencyID" validate:"required"`
		ConvertRate *uint `json:"convertRate"`
	}
	params := controller.Validated[updateCurrencyParams](ctx)
	
	currencyInfo := currencydto.UpdateCurrencyRequest{
		ID:          params.ID,
		ConvertRate: params.ConvertRate,
	}

	if err := currencyController.currencyService.UpdateCurrency(currencyInfo); err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, currencyController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateCurrency")
	controller.Response(ctx, 200, message, nil)
}