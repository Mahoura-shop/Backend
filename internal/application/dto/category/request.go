package categorydto

type CreateCategoryRequest struct {
	Name        string
	Slug        string
	Description *string
	IsActive    bool
}

type UpdateCategoryRequest struct {
	ID          uint
	Name        *string
	Slug        *string
	Description *string
	IsActive    *bool
}