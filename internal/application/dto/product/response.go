package productdto

import (
	branddto "github.com/Mahoura-shop/Backend/internal/application/dto/brand"
	categorydto "github.com/Mahoura-shop/Backend/internal/application/dto/category"
	currencydto "github.com/Mahoura-shop/Backend/internal/application/dto/currency"
)

type ProductCredential struct {
	ID            uint                            `json:"id"`
	Name          string                          `json:"name"`
	Slug          string                          `json:"slug"`
	Price         float64                         `json:"price"`
	CurrencyID    uint                            `json:"currencyID"`
	Currency      *currencydto.CurrencyCredential `json:"currency"`
	IRRPrice      uint                            `json:"irrPrice"`
	ConsumerPrice uint                            `json:"consumerPrice"`
	Step1Percent  float64                         `json:"step1Percent"`
	Step2Percent  float64                         `json:"step2Percent"`
	Step3Percent  float64                         `json:"step3Percent"`
	Step4Percent  float64                         `json:"step4Percent"`
	Step1Price    uint                            `json:"step1Price"`
	Step2Price    uint                            `json:"step2Price"`
	Step3Price    uint                            `json:"step3Price"`
	Step4Price    uint                            `json:"step4Price"`
	Step1Origin   bool                            `json:"step1Origin"`
	Step2Origin   bool                            `json:"step2Origin"`
	Step3Origin   bool                            `json:"step3Origin"`
	Step4Origin   bool                            `json:"step4Origin"`
	Quantity      uint                            `json:"quantity"`
	QuantityType  string                          `json:"quantityType"`
	Priority      uint                            `json:"priority"`
	MinOrder      uint                            `json:"minOrder"`
	Category      *categorydto.CategoryCredential `json:"category"`
	CategoryID    uint                            `json:"categoryID"`
	Brand         *branddto.BrandCredential       `json:"brand"`
	BrandID       uint                            `json:"brandID"`
	Description   string                          `json:"description"`
	Offer         string                          `json:"offer"`
	IsActive      bool                            `json:"isActive"`
	IsNew         bool                            `json:"isNew"`
	ProductPic    string                          `json:"productPic"`
	Size          string                          `json:"size"`
	Images        []string                        `json:"images"`
}

type ProductPrices struct {
	ID            uint                            `json:"id"`
	Name          string                          `json:"name"`
	Price         float64                         `json:"price"`
	IRRPrice      uint                            `json:"irrPrice"`
	Currency      *currencydto.CurrencyCredential `json:"currency"`
	NewIRRPrice   uint                            `json:"newIrrPrice"`
}