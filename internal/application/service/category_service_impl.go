package service

import (
	"errors"
	"time"

	"github.com/Mahoura-shop/Backend/bootstrap"
	categorydto "github.com/Mahoura-shop/Backend/internal/application/dto/category"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/domain/s3"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type CategoryService struct {
	constants          *bootstrap.Constants
	categoryRepository postgres.CategoryRepository
	s3Storage          s3.S3Storage
	db                 database.Database
}

type CategoryServiceDeps struct {
	Constants          *bootstrap.Constants
	CategoryRepository postgres.CategoryRepository
	S3Storage          s3.S3Storage
	DB                 database.Database
}

func NewCategoryService(deps CategoryServiceDeps) *CategoryService {
	return &CategoryService{
		constants:          deps.Constants,
		categoryRepository: deps.CategoryRepository,
		s3Storage:          deps.S3Storage,
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
		count = 0
	}
	response := categorydto.CategoryCredential{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		IsActive:    category.IsActive,
		CategoryPic: category.CategoryPic,
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
		conflictErrors.Add(categoryService.constants.Field.Slug, categoryService.constants.Tag.AlreadyExist)
		return conflictErrors
	}

	return nil
}

func (categoryService *CategoryService) GetCategories() ([]categorydto.CategoryCredential, error) {
	categories, err := categoryService.categoryRepository.GetCategories(categoryService.db)
	if err != nil {
		return nil, err
	}
	var responses []categorydto.CategoryCredential
	for _, category := range categories {
		response, err := categoryService.ParseCategory(*category)
		if err != nil {
			return nil, err
		}
		if category.CategoryPic != "" {
			categoryPic, err := categoryService.s3Storage.GetPresignedURL(enum.CategoryPic, category.CategoryPic, 8*time.Hour)
			if err != nil {
				return nil, err
			}
			response.CategoryPic = categoryPic
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (categoryService *CategoryService) CreateCategory(categoryInfo categorydto.CreateCategoryRequest) error {
	err := categoryService.validateDuplicateCategory(categoryInfo.Slug)
	if err != nil {
		return err
	}
	category := &entity.Category{
		Name:     categoryInfo.Name,
		Slug:     categoryInfo.Slug,
		IsActive: categoryInfo.IsActive,
	}

	err = categoryService.db.WithTransaction(func(tx database.Database) error {
		if categoryInfo.Description != nil {
			category.Description = *categoryInfo.Description
		} else {
			category.Description = ""
		}
		createdCategory, err := categoryService.categoryRepository.CreateCategory(tx, *category)
		if err != nil {
			return err
		}
		category.ID = createdCategory.ID
		if categoryInfo.CategoryPic != nil {
			category.CategoryPic = categoryService.constants.S3BucketPath.GetCategoryPicPath(createdCategory.ID, categoryInfo.CategoryPic.Filename)
			if err := categoryService.s3Storage.UploadObject(enum.CategoryPic, category.CategoryPic, categoryInfo.CategoryPic); err != nil {
				networkErr := exception.ClassifyNetworkError(err, "ArvanStorage", "UploadObject")
				_ = categoryService.categoryRepository.DeleteCategoryByID(tx, createdCategory.ID);
				return networkErr
			}
		
			if err := categoryService.categoryRepository.UpdateCategory(tx, *category); err != nil {
				return err
			}
		}
		return nil
	})

	return err
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
		if categoryInfo.CategoryPic != nil {
			categoryPicPath := categoryService.constants.S3BucketPath.GetCategoryPicPath(categoryInfo.ID, categoryInfo.CategoryPic.Filename)
			if err := categoryService.s3Storage.UploadObject(enum.CategoryPic, categoryPicPath, categoryInfo.CategoryPic); err != nil {
				return exception.ClassifyNetworkError(err, "ArvanStorage", "UploadObject")
			}
			if category.CategoryPic != "" {
				if err := categoryService.s3Storage.DeleteObject(enum.CategoryPic, category.CategoryPic); err != nil {
					return exception.ClassifyNetworkError(err, "ArvanStorage", "DeleteObject")
				}
			}
			category.CategoryPic = categoryPicPath
		}
		if err := categoryService.categoryRepository.UpdateCategory(tx, *category); err != nil {
			return err
		}
		return nil
	})

	return err
}