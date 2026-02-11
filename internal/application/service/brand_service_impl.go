package service

import (
	"errors"
	"time"

	"github.com/Mahoura-shop/Backend/bootstrap"
	branddto "github.com/Mahoura-shop/Backend/internal/application/dto/brand"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/domain/s3"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type BrandService struct {
	constants       *bootstrap.Constants
	brandRepository postgres.BrandRepository
	s3Storage       s3.S3Storage
	db              database.Database
}

type BrandServiceDeps struct {
	Constants       *bootstrap.Constants
	BrandRepository postgres.BrandRepository
	S3Storage       s3.S3Storage
	DB              database.Database
}

func NewBrandService(deps BrandServiceDeps) *BrandService {
	return &BrandService{
		constants:       deps.Constants,
		brandRepository: deps.BrandRepository,
		s3Storage:       deps.S3Storage,
		db:              deps.DB,
	}
}


func (brandService *BrandService) GetBrandProductsCount(brandID uint) (uint, error) {
	count, err := brandService.brandRepository.GetBrandProductsCount(brandService.db, brandID)
	if (err != nil) {
		return 0, err
	}
	return count, nil
}

func (brandService *BrandService) ParseBrand(brand entity.Brand) (branddto.BrandCredential, error) {
	count, err := brandService.GetBrandProductsCount(brand.ID)
	if err != nil {
		count = 0
	}
	response := branddto.BrandCredential{
		ID:          brand.ID,
		Name:        brand.Name,
		Slug:        brand.Slug,
		Description: brand.Description,
		IsActive:    brand.IsActive,
		BrandPic:    brand.BrandPic,
		Count:       count,
	}
	return response, nil
}


func (brandService *BrandService) FindBrandByID(brandID uint) (*branddto.BrandCredential, error) {
	brand, err := brandService.brandRepository.FindBrandByID(brandService.db, brandID)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		notFoundError := exception.NotFoundError{Item: brandService.constants.Field.Brand}
		return nil, notFoundError
	}

	parsedBrand, err := brandService.ParseBrand(*brand)
	if err != nil {
		return nil, err
	}
	return &parsedBrand, nil
}


func (brandService *BrandService) FindBrandBySlug(slug string) (*branddto.BrandCredential, error) {
	brand, err := brandService.brandRepository.FindBrandBySlug(brandService.db, slug)
	if err != nil {
		return nil, err
	}
	if brand == nil {
		notFoundError := exception.NotFoundError{Item: brandService.constants.Field.Brand}
		return nil, notFoundError
	}
	
	parsedBrand, err := brandService.ParseBrand(*brand)
	if err != nil {
		return nil, err
	}
	return &parsedBrand, nil
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
		conflictErrors.Add(brandService.constants.Field.Slug, brandService.constants.Tag.AlreadyExist)
		return conflictErrors
	}

	return nil
}

func (brandService *BrandService) GetBrands() ([]branddto.BrandCredential, error) {
	brands, err := brandService.brandRepository.GetBrands(brandService.db)
	if err != nil {
		return nil, err
	}
	var responses []branddto.BrandCredential
	for _, brand := range brands {
		response, err := brandService.ParseBrand(*brand)
		if err != nil {
			return nil, err
		}
		
		if brand.BrandPic != "" {
			brandPic, err := brandService.s3Storage.GetPresignedURL(enum.BrandPic, brand.BrandPic, 8*time.Hour)
			if err != nil {
				return nil, err
			}
			response.BrandPic = brandPic
		}
		
		responses = append(responses, response)
	}
	return responses, nil
}

func (brandService *BrandService) CreateBrand(brandInfo branddto.CreateBrandRequest) error {
	err := brandService.validateDuplicateBrand(brandInfo.Slug)
	if err != nil {
		return err
	}
	brand := &entity.Brand{
		Name:     brandInfo.Name,
		Slug:     brandInfo.Slug,
		IsActive: brandInfo.IsActive,
	}

	err = brandService.db.WithTransaction(func(tx database.Database) error {
		if brandInfo.Description != nil {
			brand.Description = *brandInfo.Description
		} else {
			brand.Description = ""
		}
		createdBrand, err := brandService.brandRepository.CreateBrand(tx, brand)
		if err != nil {
			return err
		}
		if brandInfo.BrandPic != nil {
			brand.BrandPic = brandService.constants.S3BucketPath.GetBrandPicPath(createdBrand.ID, brandInfo.BrandPic.Filename)
			if err := brandService.s3Storage.UploadObject(enum.BrandPic, brand.BrandPic, brandInfo.BrandPic); err != nil {
				return  err
			}
		
			if err := brandService.brandRepository.UpdateBrand(tx, brand); err != nil {
				return err
			}
		}
		return nil
	})

	return err
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
		if brandInfo.BrandPic != nil {
			brandPicPath := brandService.constants.S3BucketPath.GetBrandPicPath(brandInfo.ID, brandInfo.BrandPic.Filename)
			brandService.s3Storage.UploadObject(enum.BrandPic, brandPicPath, brandInfo.BrandPic)
			brand.BrandPic = brandPicPath
		}
		if err := brandService.brandRepository.UpdateBrand(tx, brand); err != nil {
			return err
		}
		return nil
	})

	return err
}