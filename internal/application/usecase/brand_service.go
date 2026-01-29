package usecase

import branddto "github.com/Mahoura-shop/Backend/internal/application/dto/brand"

type BrandService interface {
	CreateBrand(branddto.CreateBrandRequest) error
	GetBrands() ([]branddto.BrandCredential, error)
	DeleteBrand(uint) error
	UpdateBrand(branddto.UpdateBrandRequest) error
}
