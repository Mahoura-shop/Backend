package service

import (
	"errors"

	"github.com/Mahoura-shop/Backend/bootstrap"
	categorydto "github.com/Mahoura-shop/Backend/internal/application/dto/category"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type CategoryService struct {
	constants          *bootstrap.Constants
	categoryRepository postgres.CategoryRepository
	db                 database.Database
}

type CategoryServiceDeps struct {
	Constants           *bootstrap.Constants
	CategoryRepository postgres.CategoryRepository
	DB                  database.Database
}

func NewCategoryService(deps CategoryServiceDeps) *CategoryService {
	return &CategoryService{
		constants:          deps.Constants,
		categoryRepository: deps.CategoryRepository,
		db:                 deps.DB,
	}
}

func (categoryService *CategoryService) GetCategoryProductsCount(categoryID uint) (uint, error) {
	count, err := categoryService.categoryRepository.GetCategoryProductsCount(categoryService.db, categoryID)
	if (err != nil) {
		return 0, err
	}
	return count, nil
}


func (categoryService *CategoryService) ParseCategory(category entity.Category) (categorydto.CategoryCredential, error) {
	count, err := categoryService.GetCategoryProductsCount(category.ID)
	if err != nil {
		return categorydto.CategoryCredential{}, err
	}
	response := categorydto.CategoryCredential{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		IsActive:    category.IsActive,
		Count:       count,
	}
	return response, nil
}

func (categoryService *CategoryService) FindCategoryByID(categoryID uint) (*categorydto.CategoryCredential, error) {
	category, err := categoryService.categoryRepository.FindCategoryByID(categoryService.db, categoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		notFoundError := exception.NotFoundError{Item: categoryService.constants.Field.Category}
		return nil, notFoundError
	}

	parsedCategory, err := categoryService.ParseCategory(*category)
	if err != nil {
		return nil, err
	}
	return &parsedCategory, nil
}

func (categoryService *CategoryService) FindCategoryBySlug(slug string) (*categorydto.CategoryCredential, error) {
	category, err := categoryService.categoryRepository.FindCategoryBySlug(categoryService.db, slug)
	if err != nil {
		return nil, err
	}
	if category == nil {
		notFoundError := exception.NotFoundError{Item: categoryService.constants.Field.Category}
		return nil, notFoundError
	}

	parsedCategory, err := categoryService.ParseCategory(*category)
	if err != nil {
		return nil, err
	}
	return &parsedCategory, nil
}

func (categoryService *CategoryService) validateDuplicateCategory(slug string) error {
	var conflictErrors exception.ConflictErrors
	category, err := categoryService.categoryRepository.FindCategoryBySlug(categoryService.db, slug)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	
	if category != nil {
		conflictErrors.Add(categoryService.constants.Field.Category, categoryService.constants.Tag.AlreadyExist)
		return conflictErrors
	}

	return nil
}

func (categoryService *CategoryService) CreateCategory(categoryInfo categorydto.CreateCategoryRequest) error {
	err := categoryService.validateDuplicateCategory(categoryInfo.Slug)
	if err != nil {
		return err
	}

	err = categoryService.db.WithTransaction(func(tx database.Database) error {
		category := &entity.Category{
			Name:     categoryInfo.Name,
			Slug:     categoryInfo.Slug,
			IsActive: categoryInfo.IsActive,
		}
		
		if categoryInfo.Description != nil {
			category.Description = *categoryInfo.Description
		} else {
			category.Description = ""
		}
		err = categoryService.categoryRepository.CreateCategory(tx, category)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (categoryService *CategoryService) GetCategories() ([]categorydto.CategoryCredential, error) {
	categories, err := categoryService.categoryRepository.GetCategories(categoryService.db)
	if err != nil {
		return nil, err
	}
	var responses []categorydto.CategoryCredential
	for _, category := range categories {
		response := categorydto.CategoryCredential{
			ID:          category.ID,
			Name:        category.Name,
			Slug:        category.Slug,
			Description: category.Description,
			IsActive:    category.IsActive,
		}
		
		responses = append(responses, response)
	}
	return responses, nil
}

func (categoryService *CategoryService) DeleteCategory(categoryID uint) error {
	category, err := categoryService.categoryRepository.FindCategoryByID(categoryService.db, categoryID)
	if err != nil {
		return err
	}
	if category == nil {
		return exception.NotFoundError{Item: categoryService.constants.Field.Category}
	}

	if err := categoryService.categoryRepository.DeleteCategoryByID(categoryService.db, categoryID); err != nil {
		return err
	}
	return nil
}

func (categoryService *CategoryService) applyCategoryUpdates(category *entity.Category, name *string, slug *string, description *string, isActive *bool) {
	if name != nil {
		category.Name = *name
	}

	if slug != nil {
		category.Slug = *slug
	}

	if description != nil {
		category.Description = *description
	}
	
	if isActive != nil {
		category.IsActive = *isActive
	}
}

func (categoryService *CategoryService) UpdateCategory(categoryInfo categorydto.UpdateCategoryRequest) error {
	category, err := categoryService.categoryRepository.FindCategoryByID(categoryService.db, categoryInfo.ID)
	if err != nil {
		return err
	}
	if category == nil {
		return exception.NotFoundError{Item: categoryService.constants.Field.Category}
	}

	categoryService.applyCategoryUpdates(category, categoryInfo.Name, categoryInfo.Slug, categoryInfo.Description, categoryInfo.IsActive)

	if (categoryInfo.Slug != nil) {
		err := categoryService.validateDuplicateCategory(category.Slug)
		if err != nil {
			return err
		}
	}
	err = categoryService.db.WithTransaction(func(tx database.Database) error {
		if err := categoryService.categoryRepository.UpdateCategory(tx, category); err != nil {
			return err
		}
		return nil
	})

	return err
}