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
	Step4Percent  *float64
	Step1Price    *uint
	Step2Price    *uint
	Step3Price    *uint
	Step4Price    *uint
	Step1Origin   *bool
	Step2Origin   *bool
	Step3Origin   *bool
	Step4Origin   *bool
	Quantity      *uint      
	QuantityType  *string    
	Priority      *uint      
	MinOrder      *uint      
	CategoryID    *uint 
	BrandID       *uint 
	Description   *string
	Offer         *string
	IsActive      *bool
	IsNew         *bool
	ProductPic    *multipart.FileHeader 
	Size          *string
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
	Step4Percent  *float64
	Step1Price    *uint
	Step2Price    *uint
	Step3Price    *uint
	Step4Price    *uint
	Step1Origin   *bool
	Step2Origin   *bool
	Step3Origin   *bool
	Step4Origin   *bool
	Quantity      *uint      
	QuantityType  *string    
	Priority      *uint      
	MinOrder      *uint      
	CategoryID    *uint 
	BrandID       *uint 
	Description   *string
	Offer         *string
	IsActive      *bool
	IsNew         *bool
	ProductPic    *multipart.FileHeader 
	Size          *string
}

type ProductPriceUpdateCredentials struct {
	ID       uint
	IRRPrice uint
}