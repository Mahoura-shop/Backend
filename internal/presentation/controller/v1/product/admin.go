package product

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminProductController struct {
	constants       *bootstrap.Constants
	pagination      *bootstrap.Pagination
	productService usecase.ProductService
}

func NewAdminProductController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	productService usecase.ProductService,
) *AdminProductController {
	return &AdminProductController{
		constants:   constants,
		pagination:  pagination,
		productService: productService,
	}
}

func (productController *AdminProductController) CreateProduct(ctx *gin.Context) {
	type createProductParams struct {
		Name         string  `json:"name" validate:"required"`
		Slug         string  `json:"slug" validate:"required"`
		Price        float64 `json:"price" validate:"required"`
		Description  *string `json:"description"`
		IsActive     *bool   `json:"isActive"`
		IsNew        *bool   `json:"isNew"`
		Priority     *uint   `json:"priority"`
		MinOrder     *uint   `json:"minOrder"`
		CategoryID   *uint   `json:"categoryID"`
		Quantity     *uint   `json:"quantity"`
		QuantityType *string `json:"quantityType"`
		CurrencyCode *string `json:"currencyCode"`
		
	}
	params := controller.Validated[createProductParams](ctx)

	productInfo := productdto.CreateProductRequest{
		Name:         params.Name,
		Slug:         params.Slug,
		Description:  params.Description,
		IsActive:     params.IsActive,
		IsNew:        params.IsNew,
		Priority:     params.Priority,
		MinOrder:     params.MinOrder,
		CategoryID:   params.CategoryID,
		Quantity:     params.Quantity,
		QuantityType: params.QuantityType,
		Price:        params.Price,
		CurrencyCode: params.CurrencyCode,
	}
	if err := productController.productService.CreateProduct(productInfo); err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, productController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createProduct")
	controller.Response(ctx, 200, message, nil)
}

func (productController *AdminProductController) GetProducts(ctx *gin.Context) {
	products, err := productController.productService.GetProducts();
	if err != nil {
		panic(err)
	}
	
	controller.Response(ctx, 200, "", products)
}