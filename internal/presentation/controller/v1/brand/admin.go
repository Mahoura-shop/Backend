package brand

import (
	"mime/multipart"

	"github.com/Mahoura-shop/Backend/bootstrap"
	branddto "github.com/Mahoura-shop/Backend/internal/application/dto/brand"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminBrandController struct {
	constants       *bootstrap.Constants
	pagination      *bootstrap.Pagination
	brandService usecase.BrandService
}

func NewAdminBrandController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	brandService usecase.BrandService,
) *AdminBrandController {
	return &AdminBrandController{
		constants:   constants,
		pagination:  pagination,
		brandService: brandService,
	}
}

func (brandController *AdminBrandController) GetBrands(ctx *gin.Context) {
	brands, err := brandController.brandService.GetBrands(); 
	if (err != nil) {
		panic(err)
	}
	controller.Response(ctx, 200, "", brands)
}

func (brandController *AdminBrandController) CreateBrand(ctx *gin.Context) {
	type createBrandParams struct {
		Name        string                `form:"name" validate:"required"`
		Slug        string                `form:"slug" validate:"required"`
		Description *string               `form:"description"`
		IsActive    *bool                 `form:"isActive"`
		BrandPic    *multipart.FileHeader `form:"brandPic"`
	}
	params := controller.Validated[createBrandParams](ctx)

	isActive := true
	if params.IsActive != nil {
		isActive = *params.IsActive
	}

	brandInfo := branddto.CreateBrandRequest{
		Name:        params.Name,
		Slug:        params.Slug,
		Description: params.Description,
		IsActive:    isActive,
		BrandPic:    params.BrandPic,
	}
	if err := brandController.brandService.CreateBrand(brandInfo); err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, brandController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createBrand")
	controller.Response(ctx, 200, message, nil)
}

func (brandController *AdminBrandController) DeleteBrand(ctx *gin.Context) {
	type deleteBrandParams struct {
		BrandID uint `uri:"brandID" validate:"required"`
	}
	params := controller.Validated[deleteBrandParams](ctx)
	
	if err := brandController.brandService.DeleteBrand(params.BrandID); err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, brandController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteBrand")
	controller.Response(ctx, 200, message, nil)
}

func (brandController *AdminBrandController) UpdateBrand(ctx *gin.Context) {
	type updateBrandParams struct {
		ID          uint                  `uri:"brandID" validate:"required"`
		Name        *string               `form:"name"`
		Slug        *string               `form:"slug"`
		Description *string               `form:"description"`
		IsActive    *bool                 `form:"isActive"`
		BrandPic    *multipart.FileHeader `form:"brandPic"`
	}
	params := controller.Validated[updateBrandParams](ctx)
	
	brandInfo := branddto.UpdateBrandRequest{
		ID:          params.ID,
		Name:        params.Name,
		Slug:        params.Slug,
		Description: params.Description,
		IsActive:    params.IsActive,
		BrandPic:    params.BrandPic,
	}

	if err := brandController.brandService.UpdateBrand(brandInfo); err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, brandController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateBrand")
	controller.Response(ctx, 200, message, nil)
}