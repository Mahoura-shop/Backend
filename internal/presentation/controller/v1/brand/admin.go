package brand

import (
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

// func (brandController *AdminBrandController) GetCategories(ctx *gin.Context) {
// 	categories, err := brandController.brandService.GetCategories(); 
// 	if (err != nil) {
// 		panic(err)
// 	}
// 	controller.Response(ctx, 200, "", categories)
// }

func (brandController *AdminBrandController) CreateBrand(ctx *gin.Context) {
	type createBrandParams struct {
		Name        string  `json:"name" validate:"required"`
		Slug        string  `json:"slug" validate:"required"`
		Description *string `json:"description" validate:"omitempty"`
		IsActive    *bool   `json:"isActive" validate:"omitempty"`
	}
	params := controller.Validated[createBrandParams](ctx)

	brandInfo := branddto.CreateBrandRequest{
		Name:        params.Name,
		Slug:        params.Slug,
		Description: params.Description,
		IsActive:    func() bool {
			if params.IsActive != nil {
				return *params.IsActive
			}
			return true
		}(),
	}
	if err := brandController.brandService.CreateBrand(brandInfo); err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, brandController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createBrand")
	controller.Response(ctx, 200, message, nil)
}

// func (brandController *AdminBrandController) DeleteBrand(ctx *gin.Context) {
// 	type deleteBrandParams struct {
// 		BrandID uint `uri:"brandID" validate:"required"`
// 	}
// 	params := controller.Validated[deleteBrandParams](ctx)
	
// 	if err := brandController.brandService.DeleteBrand(params.BrandID); err != nil {
// 		panic(err)
// 	}
	
// 	trans := controller.GetTranslator(ctx, brandController.constants.Context.Translator)
// 	message, _ := trans.Translate("successMessage.deleteBrand")
// 	controller.Response(ctx, 200, message, nil)
// }

// func (brandController *AdminBrandController) UpdateBrand(ctx *gin.Context) {
// 	type updateBrandParams struct {
// 		ID          uint    `uri:"brandID" validate:"required"`
// 		Name        *string `json:"name"`
// 		Slug        *string `json:"slug"`
// 		Description *string `json:"description"`
// 		IsActive    *bool   `json:"isActive"`
// 	}
// 	params := controller.Validated[updateBrandParams](ctx)
	
// 	brandInfo := branddto.UpdateBrandRequest{
// 		ID:          params.ID,
// 		Name:        params.Name,
// 		Slug:        params.Slug,
// 		Description: params.Description,
// 		IsActive:    params.IsActive,
// 	}

// 	if err := brandController.brandService.UpdateBrand(brandInfo); err != nil {
// 		panic(err)
// 	}
	
// 	trans := controller.GetTranslator(ctx, brandController.constants.Context.Translator)
// 	message, _ := trans.Translate("successMessage.updateBrand")
// 	controller.Response(ctx, 200, message, nil)
// }