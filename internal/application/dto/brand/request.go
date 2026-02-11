package branddto

import "mime/multipart"

type CreateBrandRequest struct {
	Name        string
	Slug        string
	Description *string
	IsActive    bool
	BrandPic    *multipart.FileHeader
}

type UpdateBrandRequest struct {
	ID          uint
	Name        *string
	Slug        *string
	Description *string
	IsActive    *bool
	BrandPic    *multipart.FileHeader
}