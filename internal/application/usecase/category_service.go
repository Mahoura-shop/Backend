package usecase

import categorydto "github.com/Mahoura-shop/Backend/internal/application/dto/category"

type CategoryService interface {
	CreateCategory(categoryInfoInfo categorydto.CreateCategoryRequest) error
}
