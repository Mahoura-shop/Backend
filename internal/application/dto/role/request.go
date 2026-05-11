package roledto

type CreateRoleRequest struct {
	Name          string `json:"name" validate:"required"`
	Description   string `json:"description"`
	PermissionIDs []uint `json:"permissionIDs"`
}

type UpdateRoleRequest struct {
	ID            uint   `uri:"roleID" validate:"required"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	PermissionIDs []uint `json:"permissionIDs"`
}
