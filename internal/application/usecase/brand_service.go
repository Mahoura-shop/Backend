package usecase

import (
	branddto "github.com/Mahoura-shop/Backend/internal/application/dto/brand"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
)

type BrandService interface {
	ParseBrand(entity.Brand) (branddto.BrandCredential, error)
	FindBrandBySlug(string) (*branddto.BrandCredential, error) 
	FindBrandByID(uint) (*branddto.BrandCredential, error) 
	CreateBrand(branddto.CreateBrandRequest) error
	GetBrands() ([]branddto.BrandCredential, error)
	DeleteBrand(uint) error
	UpdateBrand(branddto.UpdateBrandRequest) error
}
