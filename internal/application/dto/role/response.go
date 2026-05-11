package roledto

type PermissionCredential struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RoleCredential struct {
	ID          uint                   `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Permissions []PermissionCredential `json:"permissions"`
}
