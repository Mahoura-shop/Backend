package service

import (
	"errors"

	"github.com/Mahoura-shop/Backend/bootstrap"
	branddto "github.com/Mahoura-shop/Backend/internal/application/dto/brand"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type BrandService struct {
	constants       *bootstrap.Constants
	brandRepository postgres.BrandRepository
	db              database.Database
}

type BrandServiceDeps struct {
	Constants       *bootstrap.Constants
	BrandRepository postgres.BrandRepository
	DB              database.Database
}

func NewBrandService(deps BrandServiceDeps) *BrandService {
	return &BrandService{
		constants:       deps.Constants,
		brandRepository: deps.BrandRepository,
		db:              deps.DB,
	}
}

func (brandService *BrandService) FindBrandBySlug(slug string) (*entity.Brand, error) {
	brand, err := brandService.brandRepository.FindBrandBySlug(brandService.db, slug)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		notFoundError := exception.NotFoundError{Item: brandService.constants.Field.Brand}
		return nil, notFoundError
	}
	return brand, nil
}

func (brandService *BrandService) validateDuplicateBrand(slug string) error {
	var conflictErrors exception.ConflictErrors
	brand, err := brandService.brandRepository.FindBrandBySlug(brandService.db, slug)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	if brand != nil {
		conflictErrors.Add(brandService.constants.Field.Brand, brandService.constants.Tag.AlreadyExist)
		return conflictErrors
	}

	return nil
}

func (brandService *BrandService) CreateBrand(brandInfo branddto.CreateBrandRequest) error {
	err := brandService.validateDuplicateBrand(brandInfo.Slug)
	if err != nil {
		return err
	}

	err = brandService.db.WithTransaction(func(tx database.Database) error {
		brand := &entity.Brand{
			Name:     brandInfo.Name,
			Slug:     brandInfo.Slug,
			IsActive: brandInfo.IsActive,
		}
		
		if brandInfo.Description != nil {
			brand.Description = *brandInfo.Description
		} else {
			brand.Description = ""
		}
		err = brandService.brandRepository.CreateBrand(tx, brand)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}

func (brandService *BrandService) GetBrands() ([]branddto.BrandCredential, error) {
	brands, err := brandService.brandRepository.GetBrands(brandService.db)
	if err != nil {
		return nil, err
	}
	var responses []branddto.BrandCredential
	for _, brand := range brands {
		response := branddto.BrandCredential{
			ID:       brand.ID,
			Name:     brand.Name,
			Slug:     brand.Slug,
			IsActive: brand.IsActive,
		}
		
		if brand.Description != "" {
			description := brand.Description
			response.Description = &description
		}
		
		responses = append(responses, response)
	}
	return responses, nil
}

func (brandService *BrandService) DeleteBrand(brandID uint) error {
	brand, err := brandService.brandRepository.FindBrandByID(brandService.db, brandID)
	if err != nil {
		return err
	}
	if brand == nil {
		return exception.NotFoundError{Item: brandService.constants.Field.Brand}
	}

	if err := brandService.brandRepository.DeleteBrandByID(brandService.db, brandID); err != nil {
		return err
	}
	return nil
}

func (brandService *BrandService) applyBrandUpdates(brand *entity.Brand, name *string, slug *string, description *string, isActive *bool) {
	if name != nil {
		brand.Name = *name
	}

	if slug != nil {
		brand.Slug = *slug
	}

	if description != nil {
		brand.Description = *description
	}
	
	if isActive != nil {
		brand.IsActive = *isActive
	}
}

func (brandService *BrandService) UpdateBrand(brandInfo branddto.UpdateBrandRequest) error {
	brand, err := brandService.brandRepository.FindBrandByID(brandService.db, brandInfo.ID)
	if err != nil {
		return err
	}
	if brand == nil {
		return exception.NotFoundError{Item: brandService.constants.Field.Brand}
	}

	brandService.applyBrandUpdates(brand, brandInfo.Name, brandInfo.Slug, brandInfo.Description, brandInfo.IsActive)

	if (brandInfo.Slug != nil) {
		err := brandService.validateDuplicateBrand(brand.Slug)
		if err != nil {
			return err
		}
	}
	err = brandService.db.WithTransaction(func(tx database.Database) error {
		if err := brandService.brandRepository.UpdateBrand(tx, brand); err != nil {
			return err
		}
		return nil
	})

	return err
}