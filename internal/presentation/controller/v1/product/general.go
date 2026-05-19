package product

import (
	"strconv"

	"github.com/Mahoura-shop/Backend/bootstrap"
	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	domainPostgres "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type GeneralProductController struct {
	constants       *bootstrap.Constants
	productService  usecase.ProductService
	userRepository  domainPostgres.UserRepository
	database        database.Database
}

func NewGeneralProductController(
	constants *bootstrap.Constants,
	productService usecase.ProductService,
	userRepository domainPostgres.UserRepository,
	database database.Database,
) *GeneralProductController {
	return &GeneralProductController{
		constants:       constants,
		productService:  productService,
		userRepository:  userRepository,
		database:        database,
	}
}

func (c *GeneralProductController) GetProducts(ctx *gin.Context) {
	var filter productdto.ProductFilterRequest
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		panic(err)
	}

	result, err := c.productService.SearchProductsWithPagination(filter)
	if err != nil {
		panic(err)
	}

	var userType enum.UserType
	userID, exists := ctx.Get(c.constants.Context.ID)
	if exists && userID != nil {
		user, err := c.userRepository.FindUserByID(c.database, userID.(uint))
		if err == nil && user != nil {
			userType = user.Type
		}
	}

	for i := range result.Products {
		result.Products[i].ResolvedPrice = c.resolvePriceForType(userType, &result.Products[i])
	}

	controller.Response(ctx, 200, "", result)
}

func firstNonZero(vals ...uint) uint {
	for _, v := range vals {
		if v != 0 {
			return v
		}
	}
	return 0
}

func (c *GeneralProductController) resolvePriceForType(userType enum.UserType, product *productdto.ProductCredential) uint {
	switch userType {
	case enum.UserTypeFellow, enum.UserTypeAdmin:
		return firstNonZero(product.Step1Price, product.Step2Price, product.Step3Price, product.Step4Price, product.ConsumerPrice, product.IRRPrice)
	case enum.UserTypeShopkeeperCash, enum.UserTypeShopkeeperCheque:
		return firstNonZero(product.Step2Price, product.Step3Price, product.Step4Price, product.ConsumerPrice, product.IRRPrice)
	case enum.UserTypeCustomer:
		return firstNonZero(product.Step4Price, product.ConsumerPrice, product.IRRPrice)
	default:
		return firstNonZero(product.ConsumerPrice, product.IRRPrice)
	}
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
