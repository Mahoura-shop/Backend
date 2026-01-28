package categorydto

type CreateCategoryRequest struct {
	Name        string
	Slug        string
	Description *string
	IsActive    bool
}