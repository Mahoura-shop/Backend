package productdto

type ProductCredential struct {
	ID           uint    `json:"id"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Price        float64 `json:"price"` 
	Description  string  `json:"description"`
	IsActive     bool    `json:"isActive"`
	IsNew        bool    `json:"isNew"`
	Priority     uint    `json:"priority"`
	MinOrder     uint    `json:"minOrder"`
	CategoryID   *uint   `json:"categoryId"`
	BrandID      *uint   `json:"brandId"`
	Quantity     uint    `json:"quantity"`
	QuantityType string  `json:"quantityType"`  
	CurrencyCode string  `json:"currencyCode"`
	ProductPic   string  `json:"productPic"`
}