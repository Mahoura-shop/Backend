package user

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

const adminTargetUserID = uint(99)

type AdminUserTestSuite struct {
	suite.Suite
	constants   *bootstrap.Constants
	userService *mocks.UserServiceMock
	controller  *AdminUserController
}

func (s *AdminUserTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.userService = mocks.NewUserServiceMock()
	s.controller = NewAdminUserController(s.constants, &bootstrap.Pagination{}, s.userService)
}

func (s *AdminUserTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Next()
	})
	if withRecovery {
		r.Use(adminUserPanicToHTTP())
	}
	r.GET("/admin/users", s.controller.GetUsers)
	r.POST("/admin/users/:userID/ban", s.controller.BanUser)
	r.POST("/admin/users/:userID/unban", s.controller.UnbanUser)
	r.GET("/admin/dashboard", s.controller.GetDashboard)
	return r
}

func adminUserPanicToHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case exception.NotFoundError:
					c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
				case exception.ValidationErrors:
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "validation"})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

// ADM-20: GET /admin/users → 200 with all users
func (s *AdminUserTestSuite) TestADM20_GetUsers_Returns200() {
	s.userService.On("GetUsers").Return([]userdto.UserCredential{
		{ID: 1, Phone: "+989164911318", Type: "regular"},
		{ID: 2, Phone: "+989164911319", Type: "shopkeeperCash"},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin/users", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	users := body["data"].([]any)
	s.Len(users, 2)
	s.userService.AssertExpectations(s.T())
}

// ADM-20b: Empty user list → 200 with empty array
func (s *AdminUserTestSuite) TestADM20b_GetUsers_Empty_Returns200() {
	s.userService.On("GetUsers").Return([]userdto.UserCredential{}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin/users", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	users := body["data"].([]any)
	s.Len(users, 0)
	s.userService.AssertExpectations(s.T())
}

// ADM-26: POST /admin/users/:userID/ban → 200, user banned
func (s *AdminUserTestSuite) TestADM26_BanUser_Returns200() {
	s.userService.On("BanUser", adminTargetUserID).Return(nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/admin/users/99/ban", nil))

	s.Equal(http.StatusOK, w.Code)
	s.userService.AssertExpectations(s.T())
}

// ADM-26b: Ban non-existent user → service returns NotFoundError → 404
func (s *AdminUserTestSuite) TestADM26b_BanUser_NotFound_Returns404() {
	s.userService.On("BanUser", adminTargetUserID).Return(
		exception.NotFoundError{Item: "user"},
	)

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/admin/users/99/ban", nil))

	s.Equal(http.StatusNotFound, w.Code)
	s.userService.AssertExpectations(s.T())
}

// ADM-27: POST /admin/users/:userID/unban → 200, user unbanned
func (s *AdminUserTestSuite) TestADM27_UnbanUser_Returns200() {
	s.userService.On("UnbanUser", adminTargetUserID).Return(nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/admin/users/99/unban", nil))

	s.Equal(http.StatusOK, w.Code)
	s.userService.AssertExpectations(s.T())
}

// ADM-27b: Unban non-existent user → service returns NotFoundError → 404
func (s *AdminUserTestSuite) TestADM27b_UnbanUser_NotFound_Returns404() {
	s.userService.On("UnbanUser", adminTargetUserID).Return(
		exception.NotFoundError{Item: "user"},
	)

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/admin/users/99/unban", nil))

	s.Equal(http.StatusNotFound, w.Code)
	s.userService.AssertExpectations(s.T())
}

// ADM-70: GET /admin/dashboard → 200 with revenue, order counts, user counts
func (s *AdminUserTestSuite) TestADM70_GetDashboard_Returns200() {
	s.userService.On("GetDashboard").Return(userdto.DashboardResponse{
		ProductsCount:   50,
		CategoriesCount: 10,
		BrandsCount:     8,
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	s.Equal(float64(50), data["productsCount"])
	s.userService.AssertExpectations(s.T())
}

func TestAdminUserSuite(t *testing.T) {
	suite.Run(t, new(AdminUserTestSuite))
}
