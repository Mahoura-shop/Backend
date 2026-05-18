package role

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	roledto "github.com/Mahoura-shop/Backend/internal/application/dto/role"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

func init() { gin.SetMode(gin.TestMode) }

const testRoleID = uint(2)

type AdminRoleTestSuite struct {
	suite.Suite
	constants   *bootstrap.Constants
	roleService *mocks.RoleServiceMock
	controller  *AdminRoleController
}

func (s *AdminRoleTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.roleService = mocks.NewRoleServiceMock()
	s.controller = NewAdminRoleController(s.constants, s.roleService)
}

func (s *AdminRoleTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Next()
	})
	if withRecovery {
		r.Use(rolePanicToHTTP())
	}
	r.GET("/admin/roles", s.controller.GetRoles)
	r.GET("/admin/permissions", s.controller.GetPermissions)
	r.POST("/admin/roles", s.controller.CreateRole)
	r.PATCH("/admin/roles/:roleID", s.controller.UpdateRole)
	r.DELETE("/admin/roles/:roleID", s.controller.DeleteRole)
	return r
}

func rolePanicToHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case exception.NotFoundError:
					c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
				case exception.ValidationErrors:
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "validation"})
				case exception.ConflictErrors:
					c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

// ADM-50: GET /admin/roles → 200 with all roles and their permissions
func (s *AdminRoleTestSuite) TestADM50_GetRoles_Returns200() {
	s.roleService.On("GetRoles").Return([]roledto.RoleCredential{
		{ID: 1, Name: "manager", Permissions: []roledto.PermissionCredential{{ID: 1, Name: "view_orders"}}},
		{ID: testRoleID, Name: "support", Permissions: []roledto.PermissionCredential{}},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin/roles", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	roles := body["data"].([]any)
	s.Len(roles, 2)
	s.roleService.AssertExpectations(s.T())
}

// ADM-54: GET /admin/permissions → 200 with full permission list
func (s *AdminRoleTestSuite) TestADM54_GetPermissions_Returns200() {
	s.roleService.On("GetPermissions").Return([]roledto.PermissionCredential{
		{ID: 1, Name: "view_orders"},
		{ID: 2, Name: "manage_users"},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin/permissions", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	perms := body["data"].([]any)
	s.Len(perms, 2)
	s.roleService.AssertExpectations(s.T())
}

// ADM-51: POST /admin/roles → 201 with created role
func (s *AdminRoleTestSuite) TestADM51_CreateRole_Returns201() {
	req := roledto.CreateRoleRequest{
		Name:          "editor",
		Description:   "can edit content",
		PermissionIDs: []uint{1, 2},
	}
	s.roleService.On("CreateRole", req).Return(&roledto.RoleCredential{
		ID:   3,
		Name: "editor",
	}, nil)

	payload, _ := json.Marshal(map[string]any{
		"name":          "editor",
		"description":   "can edit content",
		"permissionIDs": []uint{1, 2},
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/roles", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusCreated, w.Code)
	s.roleService.AssertExpectations(s.T())
}

// ADM-51b: POST /admin/roles missing required name → 422
func (s *AdminRoleTestSuite) TestADM51b_CreateRole_MissingName_Returns422() {
	payload, _ := json.Marshal(map[string]any{"description": "some role"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/roles", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(true).ServeHTTP(w, r)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
	s.roleService.AssertNotCalled(s.T(), "CreateRole")
}

// ADM-52: PATCH /admin/roles/:roleID → 200, role updated with new permissions
func (s *AdminRoleTestSuite) TestADM52_UpdateRole_Returns200() {
	req := roledto.UpdateRoleRequest{
		ID:            testRoleID,
		Name:          "support-v2",
		PermissionIDs: []uint{3},
	}
	s.roleService.On("UpdateRole", req).Return(nil)

	payload, _ := json.Marshal(map[string]any{
		"name":          "support-v2",
		"permissionIDs": []uint{3},
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/admin/roles/2", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.roleService.AssertExpectations(s.T())
}

// ADM-53: DELETE /admin/roles/:roleID → 200, role removed
func (s *AdminRoleTestSuite) TestADM53_DeleteRole_Returns200() {
	s.roleService.On("DeleteRole", testRoleID).Return(nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodDelete, "/admin/roles/2", nil))

	s.Equal(http.StatusOK, w.Code)
	s.roleService.AssertExpectations(s.T())
}

// ADM-53b: Delete non-existent role → 404
func (s *AdminRoleTestSuite) TestADM53b_DeleteRole_NotFound_Returns404() {
	s.roleService.On("DeleteRole", testRoleID).Return(exception.NotFoundError{Item: "role"})

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w,
		httptest.NewRequest(http.MethodDelete, "/admin/roles/2", nil))

	s.Equal(http.StatusNotFound, w.Code)
	s.roleService.AssertExpectations(s.T())
}

func TestAdminRoleSuite(t *testing.T) {
	suite.Run(t, new(AdminRoleTestSuite))
}
