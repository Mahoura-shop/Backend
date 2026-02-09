package categorydto

import "mime/multipart"

type CreateCategoryRequest struct {
	Name        string
	Slug        string
	Description *string
	IsActive    bool
	CategoryPic *multipart.FileHeader
}

type UpdateCategoryRequest struct {
	ID          uint
	Name        *string
	Slug        *string
	Description *string
	IsActive    *bool
	CategoryPic *multipart.FileHeader
}