package category

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	categorydto "github.com/Mahoura-shop/Backend/internal/application/dto/category"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminCategoryController struct {
	constants       *bootstrap.Constants
	pagination      *bootstrap.Pagination
	categoryService usecase.CategoryService
}

func NewAdminCategoryController(
	constants *bootstrap.Constants,
	pagination *bootstrap.Pagination,
	categoryService usecase.CategoryService,
) *AdminCategoryController {
	return &AdminCategoryController{
		constants:   constants,
		pagination:  pagination,
		categoryService: categoryService,
	}
}

func (categoryController *AdminCategoryController) GetCategories(ctx *gin.Context) {
	categories, err := categoryController.categoryService.GetCategories(); 
	if (err != nil) {
		panic(err)
	}
	controller.Response(ctx, 200, "", categories)
}

func (categoryController *AdminCategoryController) CreateCategory(ctx *gin.Context) {
	type createCategoryParams struct {
		Name        string  `json:"name" validate:"required"`
		Slug        string  `json:"slug" validate:"required"`
		Description *string `json:"description" validate:"omitempty"`
		IsActive    *bool   `json:"isActive" validate:"omitempty"`
	}
	params := controller.Validated[createCategoryParams](ctx)

	categoryInfo := categorydto.CreateCategoryRequest{
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
	if err := categoryController.categoryService.CreateCategory(categoryInfo); err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, categoryController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.createCategory")
	controller.Response(ctx, 200, message, nil)
}

func (categoryController *AdminCategoryController) DeleteCategory(ctx *gin.Context) {
	type deleteCategoryParams struct {
		CategoryID uint `uri:"categoryID" validate:"required"`
	}
	params := controller.Validated[deleteCategoryParams](ctx)
	
	if err := categoryController.categoryService.DeleteCategory(params.CategoryID); err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, categoryController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.deleteCategory")
	controller.Response(ctx, 200, message, nil)
}

func (categoryController *AdminCategoryController) UpdateCategory(ctx *gin.Context) {
	type updateCategoryParams struct {
		ID          uint    `uri:"categoryID" validate:"required"`
		Name        *string `json:"name"`
		Slug        *string `json:"slug"`
		Description *string `json:"description"`
		IsActive    *bool   `json:"isActive"`
	}
	params := controller.Validated[updateCategoryParams](ctx)
	
	categoryInfo := categorydto.UpdateCategoryRequest{
		ID:          params.ID,
		Name:        params.Name,
		Slug:        params.Slug,
		Description: params.Description,
		IsActive:    params.IsActive,
	}

	if err := categoryController.categoryService.UpdateCategory(categoryInfo); err != nil {
		panic(err)
	}
	
	trans := controller.GetTranslator(ctx, categoryController.constants.Context.Translator)
	message, _ := trans.Translate("successMessage.updateCategory")
	controller.Response(ctx, 200, message, nil)
}