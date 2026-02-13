package product

import (
	"mime/multipart"

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

func (productController *AdminProductController) GetProduct(ctx *gin.Context) {
	type getProductParams struct {
		ProductID uint `uri:"productID" validate:"required"`
	}
	params := controller.Validated[getProductParams](ctx)
	
	product, err := productController.productService.GetProduct(params.ProductID); 
	if (err != nil) {
		panic(err)
	}

	controller.Response(ctx, 200, "", product)
}

func (productController *AdminProductController) GetProducts(ctx *gin.Context) {
	products, err := productController.productService.GetProducts();
	if err != nil {
		panic(err)
	}

	controller.Response(ctx, 200, "", products)
}

func (productController *AdminProductController) CreateProduct(ctx *gin.Context) {
	type createProductParams struct {
		Name          string                `form:"name" validate:"required"`
		Slug          string                `form:"slug" validate:"required"`
		Price         float64               `form:"price" validate:"required"`
		CurrencyID    uint                  `form:"currencyID"`
		IRRPrice      *uint                 `form:"irrPrice"`
		ConsumerPrice *uint                 `form:"consumerPrice"`
		Step1Percent  *float64              `form:"step1Percent"`
		Step2Percent  *float64              `form:"step2Percent"`
		Step3Percent  *float64              `form:"step3Percent"`
		Step1Price    *uint                 `form:"step1Price"`
		Step2Price    *uint                 `form:"step2Price"`
		Step3Price    *uint                 `form:"step3Price"`
		Quantity      *uint                 `form:"quantity"`
		QuantityType  *string               `form:"quantityType"`
		Priority      *uint                 `form:"priority"`
		MinOrder      *uint                 `form:"minOrder"`
		CategoryID    *uint                 `form:"categoryID"`
		BrandID       *uint                 `form:"brandID"`
		Description   *string               `form:"description"`
		IsActive      *bool                 `form:"isActive"`
		IsNew         *bool                 `form:"isNew"`
		ProductPic    *multipart.FileHeader `form:"productPic"`
	}
	params := controller.Validated[createProductParams](ctx)

	productInfo := productdto.CreateProductRequest{
		Name:          params.Name,
		Slug:          params.Slug,
		Price:         params.Price,
		CurrencyID:    params.CurrencyID,
		IRRPrice:      params.IRRPrice,
		ConsumerPrice: params.ConsumerPrice,
		Step1Percent:  params.Step1Percent,
		Step2Percent:  params.Step2Percent,
		Step3Percent:  params.Step3Percent,
		Step1Price:    params.Step1Price,
		Step2Price:    params.Step2Price,
		Step3Price:    params.Step3Price,
		Quantity:      params.Quantity,
		QuantityType:  params.QuantityType,
		Priority:      params.Priority,
		MinOrder:      params.MinOrder,
		CategoryID:    params.CategoryID,
		BrandID:       params.BrandID,
		Description:   params.Description,
		IsActive:      params.IsActive,
		IsNew:         params.IsNew,
		ProductPic:    params.ProductPic,
	}
	if err := productController.productService.CreateProduct(productInfo); err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, productController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createProduct")
	controller.Response(ctx, 200, message, nil)
}

func (productController *AdminProductController) UpdateProduct(ctx *gin.Context) {
	type updateProductParams struct {
		ID            uint                  `uri:"productID" validate:"required"`
		Name          *string               `form:"name"`
		Slug          *string               `form:"slug"`
		Price         *float64              `form:"price"`
		CurrencyID    *uint                 `form:"currencyID"`
		IRRPrice      *uint                 `form:"irrPrice"`
		ConsumerPrice *uint                 `form:"consumerPrice"`
		Step1Percent  *float64              `form:"step1Percent"`
		Step2Percent  *float64              `form:"step2Percent"`
		Step3Percent  *float64              `form:"step3Percent"`
		Step1Price    *uint                 `form:"step1Price"`
		Step2Price    *uint                 `form:"step2Price"`
		Step3Price    *uint                 `form:"step3Price"`
		Quantity      *uint                 `form:"quantity"`
		QuantityType  *string               `form:"quantityType"`
		Priority      *uint                 `form:"priority"`
		MinOrder      *uint                 `form:"minOrder"`
		CategoryID    *uint                 `form:"categoryID"`
		BrandID       *uint                 `form:"brandID"`
		Description   *string               `form:"description"`
		IsActive      *bool                 `form:"isActive"`
		IsNew         *bool                 `form:"isNew"`
		ProductPic    *multipart.FileHeader `form:"productPic"`
	}
	params := controller.Validated[updateProductParams](ctx)
	
	productInfo := productdto.UpdateProductRequest{
		ID:            params.ID,
		Name:          params.Name,
		Slug:          params.Slug,
		Price:         params.Price,
		CurrencyID:    params.CurrencyID,
		IRRPrice:      params.IRRPrice,
		ConsumerPrice: params.ConsumerPrice,
		Step1Percent:  params.Step1Percent,
		Step2Percent:  params.Step2Percent,
		Step3Percent:  params.Step3Percent,
		Step1Price:    params.Step1Price,
		Step2Price:    params.Step2Price,
		Step3Price:    params.Step3Price,
		Quantity:      params.Quantity,
		QuantityType:  params.QuantityType,
		Priority:      params.Priority,
		MinOrder:      params.MinOrder,
		CategoryID:    params.CategoryID,
		BrandID:       params.BrandID,
		Description:   params.Description,
		IsActive:      params.IsActive,
		IsNew:         params.IsNew,
		ProductPic:    params.ProductPic,
	}

	if err := productController.productService.UpdateProduct(productInfo); err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, productController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateProduct")
	controller.Response(ctx, 200, message, nil)
}

func (productController *AdminProductController) DeleteProduct(ctx *gin.Context) {
	type deleteProductParams struct {
		ProductID uint `uri:"productID" validate:"required"`
	}
	params := controller.Validated[deleteProductParams](ctx)

	if err := productController.productService.DeleteProduct(params.ProductID); err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, productController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteProduct")
	controller.Response(ctx, 200, message, nil)
}