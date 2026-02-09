package usecase

import (
	categorydto "github.com/Mahoura-shop/Backend/internal/application/dto/category"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
)

type CategoryService interface {
	ParseCategory(entity.Category) (categorydto.CategoryCredential, error)
	FindCategoryBySlug(string) (*categorydto.CategoryCredential, error) 
	FindCategoryByID(uint) (*categorydto.CategoryCredential, error) 
	CreateCategory(categorydto.CreateCategoryRequest) error
	GetCategories() ([]categorydto.CategoryCredential, error)
	DeleteCategory(uint) error
	UpdateCategory(categorydto.UpdateCategoryRequest) error
}
