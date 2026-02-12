package productdto

import (
	branddto "github.com/Mahoura-shop/Backend/internal/application/dto/brand"
	categorydto "github.com/Mahoura-shop/Backend/internal/application/dto/category"
)

type cotrProductCredential struct {
	ID           uint                            `json:"id"`
	Name         string                          `json:"name"`
	Slug         string                          `json:"slug"`
	Price        float64                         `json:"price"`
	Description  string                          `json:"description"`
	IsActive     bool                            `json:"isActive"`
	IsNew        bool                            `json:"isNew"`
	Priority     uint                            `json:"priority"`
	MinOrder     uint                            `json:"minOrder"`
	Category     *categorydto.CategoryCredential `json:"category"`
	CategoryID   uint                            `json:"categoryID"`
	Brand        *branddto.BrandCredential       `json:"brand"`
	BrandID      uint                            `json:"brandID"`
	Quantity     uint                            `json:"quantity"`
	QuantityType string                          `json:"quantityType"`
	CurrencyCode string                          `json:"currencyCode"`
	ProductPic   string                          `json:"productPic"`
}