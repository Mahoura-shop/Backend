package categorydto

type CategoryCredentialResponse struct {
	Name        string
	Slug        string
	Description *string
	IsActive    bool
}