package productdto

import "mime/multipart"

type CreateProductRequest struct {
	Name         string
	Slug         string
	Price        float64    
	Description  *string
	IsActive     *bool
	IsNew        *bool
	Priority     *uint      
	MinOrder     *uint      
	CategoryID   *uint 
	BrandID      *uint 
	Quantity     *uint      
	QuantityType *string    
	CurrencyCode *string    
	ProductPic   *multipart.FileHeader 
}

type UpdateProductRequest struct {
	ID 	         uint
	Name         *string
	Slug         *string
	Price        *float64    
	Description  *string
	IsActive     *bool
	IsNew        *bool
	Priority     *uint      
	MinOrder     *uint      
	CategoryID   *uint 
	BrandID      *uint 
	Quantity     *uint      
	QuantityType *string    
	CurrencyCode *string     
	ProductPic   *multipart.FileHeader 
}