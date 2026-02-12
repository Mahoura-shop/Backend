package productdto

import (
	branddto "github.com/Mahoura-shop/Backend/internal/application/dto/brand"
	categorydto "github.com/Mahoura-shop/Backend/internal/application/dto/category"
)

type ProductCredential struct {
	ID            uint                            `json:"id"`
	Name          string                          `json:"name"`
	Slug          string                          `json:"slug"`
	Price         float64                         `json:"price"`
	CurrencyCode  string                          `json:"currencyCode"`
	IRRPrice      float64                         `json:"irrPrice"`
	ConsumerPrice float64                         `json:"consumerPrice"`
	Step1Percent  float64                         `json:"step1Percent"`
	Step2Percent  float64                         `json:"step2Percent"`
	Step3Percent  float64                         `json:"step3Percent"`
	Step1Price    float64                         `json:"step1Price"`
	Step2Price    float64                         `json:"step2Price"`
	Step3Price    float64                         `json:"step3Price"`
	Quantity      uint                            `json:"quantity"`
	QuantityType  string                          `json:"quantityType"`
	Priority      uint                            `json:"priority"`
	MinOrder      uint                            `json:"minOrder"`
	Category      *categorydto.CategoryCredential `json:"category"`
	CategoryID    uint                            `json:"categoryID"`
	Brand         *branddto.BrandCredential       `json:"brand"`
	BrandID       uint                            `json:"brandID"`
	Description   string                          `json:"description"`
	IsActive      bool                            `json:"isActive"`
	IsNew         bool                            `json:"isNew"`
	ProductPic    string                          `json:"productPic"`
}