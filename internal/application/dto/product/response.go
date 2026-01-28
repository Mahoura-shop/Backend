package productdto

type ProductCredential struct {
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Description  string  `json:"description"`
	IsActive     bool    `json:"isActive"`
	IsNew        bool    `json:"isNew"`
	Priority     uint    `json:"priority"`
	MinOrder     uint    `json:"minOrder"`
	CategoryID   *uint   `json:"categoryId"`
	Quantity     uint    `json:"quantity"`
	QuantityType string  `json:"quantityType"`  
	Price        float64 `json:"price"` 
	CurrencyCode string  `json:"currencyCode"`
}