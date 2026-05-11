package role

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	roledto "github.com/Mahoura-shop/Backend/internal/application/dto/role"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type AdminRoleController struct {
	constants   *bootstrap.Constants
	roleService usecase.RoleService
}

func NewAdminRoleController(
	constants *bootstrap.Constants,
	roleService usecase.RoleService,
) *AdminRoleController {
	return &AdminRoleController{
		constants:   constants,
		roleService: roleService,
	}
}

func (c *AdminRoleController) GetRoles(ctx *gin.Context) {
	roles, err := c.roleService.GetRoles()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", roles)
}

func (c *AdminRoleController) GetPermissions(ctx *gin.Context) {
	perms, err := c.roleService.GetPermissions()
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", perms)
}

func (c *AdminRoleController) CreateRole(ctx *gin.Context) {
	type params struct {
		Name          string `json:"name" validate:"required"`
		Description   string `json:"description"`
		PermissionIDs []uint `json:"permissionIDs"`
	}
	p := controller.Validated[params](ctx)

	req := roledto.CreateRoleRequest{
		Name:          p.Name,
		Description:   p.Description,
		PermissionIDs: p.PermissionIDs,
	}
	role, err := c.roleService.CreateRole(req)
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 201, "", role)
}

func (c *AdminRoleController) UpdateRole(ctx *gin.Context) {
	type params struct {
		ID            uint   `uri:"roleID" validate:"required"`
		Name          string `json:"name"`
		Description   string `json:"description"`
		PermissionIDs []uint `json:"permissionIDs"`
	}
	p := controller.Validated[params](ctx)

	req := roledto.UpdateRoleRequest{
		ID:            p.ID,
		Name:          p.Name,
		Description:   p.Description,
		PermissionIDs: p.PermissionIDs,
	}
	if err := c.roleService.UpdateRole(req); err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", nil)
}

func (c *AdminRoleController) DeleteRole(ctx *gin.Context) {
	type params struct {
		RoleID uint `uri:"roleID" validate:"required"`
	}
	p := controller.Validated[params](ctx)
	if err := c.roleService.DeleteRole(p.RoleID); err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", nil)
}
