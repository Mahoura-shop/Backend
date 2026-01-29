package branddto

type CreateBrandRequest struct {
	Name        string
	Slug        string
	Description *string
	IsActive    bool
}

type UpdateBrandRequest struct {
	ID          uint
	Name        *string
	Slug        *string
	Description *string
	IsActive    *bool
}