package product

import (
	"strconv"

	"github.com/Mahoura-shop/Backend/bootstrap"
	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralProductController struct {
	constants      *bootstrap.Constants
	productService usecase.ProductService
}

func NewGeneralProductController(
	constants *bootstrap.Constants,
	productService usecase.ProductService,
) *GeneralProductController {
	return &GeneralProductController{
		constants:      constants,
		productService: productService,
	}
}

func (c *GeneralProductController) GetProducts(ctx *gin.Context) {
	var filter productdto.ProductFilterRequest
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		panic(err)
	}
	products, err := c.productService.SearchProducts(filter)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", products)
}

func (c *GeneralProductController) GetProductBySlug(ctx *gin.Context) {
	slug := ctx.Param("slug")
	product, err := c.productService.GetProductBySlug(slug)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", product)
}

func (c *GeneralProductController) GetRelatedProducts(ctx *gin.Context) {
	productIDStr := ctx.Param("productID")
	productID, err := strconv.ParseUint(productIDStr, 10, 64)
	if err != nil {
		panic(err)
	}
	products, err := c.productService.GetRelatedProducts(uint(productID), 6)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", products)
}
