package branddto

type BrandCredential struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description *string `json:"description"`
	IsActive    bool    `json:"isActive"`
}