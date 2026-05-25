package branddto

import currencydto "github.com/Mahoura-shop/Backend/internal/application/dto/currency"

type BrandProduct struct {
	ID         uint                            `json:"id"`
	Name       string                          `json:"name"`
	Slug       string                          `json:"slug"`
	IsActive   bool                            `json:"isActive"`
	Quantity   uint                            `json:"quantity"`
	Price      float64                         `json:"price"`
	CurrencyID uint                            `json:"currencyID"`
	Currency   *currencydto.CurrencyCredential `json:"currency"`
	IRRPrice   uint                            `json:"irrPrice"`
}

type BrandCredential struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	IsActive    bool           `json:"isActive"`
	Count       uint           `json:"count"`
	BrandPic    string         `json:"brandPic"`
	Products    []BrandProduct `json:"products"`
}
