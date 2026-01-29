package usecase

import categorydto "github.com/Mahoura-shop/Backend/internal/application/dto/category"

type CategoryService interface {
	CreateCategory(categorydto.CreateCategoryRequest) error
	GetCategories() ([]categorydto.CategoryCredential, error)
	DeleteCategory(uint) error
	UpdateCategory(categorydto.UpdateCategoryRequest) error
}
