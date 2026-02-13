package productdto

import "mime/multipart"

type CreateProductRequest struct {
	Name          string
	Slug          string
	Price         float64   
	CurrencyID    uint
	IRRPrice      *uint 
	ConsumerPrice *uint
	Step1Percent  *float64
	Step2Percent  *float64
	Step3Percent  *float64
	Step1Price    *uint
	Step2Price    *uint
	Step3Price    *uint
	Quantity      *uint      
	QuantityType  *string    
	Priority      *uint      
	MinOrder      *uint      
	CategoryID    *uint 
	BrandID       *uint 
	Description   *string
	IsActive      *bool
	IsNew         *bool
	ProductPic    *multipart.FileHeader 
}

type UpdateProductRequest struct {
	ID 	          uint
	Name          *string
	Slug          *string
	Price         *float64    
	CurrencyID    *uint
	IRRPrice      *uint 
	ConsumerPrice *uint
	Step1Percent  *float64
	Step2Percent  *float64
	Step3Percent  *float64
	Step1Price    *uint
	Step2Price    *uint
	Step3Price    *uint
	Quantity      *uint      
	QuantityType  *string    
	Priority      *uint      
	MinOrder      *uint      
	CategoryID    *uint 
	BrandID       *uint 
	Description   *string
	IsActive      *bool
	IsNew         *bool
	ProductPic    *multipart.FileHeader 
}