package category

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	categorydto "github.com/Mahoura-shop/Backend/internal/application/dto/category"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

func init() { gin.SetMode(gin.TestMode) }

const catID = uint(3)

type AdminCategoryTestSuite struct {
	suite.Suite
	constants       *bootstrap.Constants
	categoryService *mocks.CategoryServiceMock
	controller      *AdminCategoryController
}

func (s *AdminCategoryTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.categoryService = mocks.NewCategoryServiceMock()
	s.controller = NewAdminCategoryController(s.constants, &bootstrap.Pagination{}, s.categoryService)
}

func (s *AdminCategoryTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Next()
	})
	if withRecovery {
		r.Use(catPanicToHTTP())
	}
	r.GET("/admin/categories", s.controller.GetCategories)
	r.POST("/admin/categories", s.controller.CreateCategory)
	r.PATCH("/admin/categories/:categoryID", s.controller.UpdateCategory)
	r.DELETE("/admin/categories/:categoryID", s.controller.DeleteCategory)
	return r
}

func catPanicToHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case exception.NotFoundError:
					c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
				case exception.ValidationErrors:
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "validation"})
				case exception.BindingError:
					c.JSON(http.StatusBadRequest, gin.H{"error": "binding"})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

func multipartBody(fields map[string]string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for k, v := range fields {
		_ = w.WriteField(k, v)
	}
	w.Close()
	return body, w.FormDataContentType()
}

// ADM-01: POST /admin/categories → 200, category saved
func (s *AdminCategoryTestSuite) TestADM01_CreateCategory_Returns200() {
	req := categorydto.CreateCategoryRequest{
		Name:     "Skin Care",
		Slug:     "skin-care",
		IsActive: true,
	}
	s.categoryService.On("CreateCategory", req).Return(nil)

	body, ct := multipartBody(map[string]string{"name": "Skin Care", "slug": "skin-care"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/categories", body)
	r.Header.Set("Content-Type", ct)
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.categoryService.AssertExpectations(s.T())
}

// ADM-01b: POST /admin/categories missing name → 422
func (s *AdminCategoryTestSuite) TestADM01b_CreateCategory_MissingName_Returns422() {
	body, ct := multipartBody(map[string]string{"slug": "skin-care"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/categories", body)
	r.Header.Set("Content-Type", ct)
	s.newRouter(true).ServeHTTP(w, r)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
	s.categoryService.AssertNotCalled(s.T(), "CreateCategory")
}

// ADM-02: PATCH /admin/categories/:categoryID → 200, category updated
func (s *AdminCategoryTestSuite) TestADM02_UpdateCategory_Returns200() {
	newName := "Skin Care Updated"
	req := categorydto.UpdateCategoryRequest{
		ID:   catID,
		Name: &newName,
	}
	s.categoryService.On("UpdateCategory", req).Return(nil)

	body, ct := multipartBody(map[string]string{"name": newName})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/admin/categories/3", body)
	r.Header.Set("Content-Type", ct)
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.categoryService.AssertExpectations(s.T())
}

// ADM-03: DELETE /admin/categories/:categoryID → 200, category removed
func (s *AdminCategoryTestSuite) TestADM03_DeleteCategory_Returns200() {
	s.categoryService.On("DeleteCategory", catID).Return(nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodDelete, "/admin/categories/3", nil))

	s.Equal(http.StatusOK, w.Code)
	s.categoryService.AssertExpectations(s.T())
}

// ADM-03b: Delete non-existent category → service returns NotFoundError → 404
func (s *AdminCategoryTestSuite) TestADM03b_DeleteCategory_NotFound_Returns404() {
	s.categoryService.On("DeleteCategory", catID).Return(exception.NotFoundError{Item: "category"})

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w,
		httptest.NewRequest(http.MethodDelete, "/admin/categories/3", nil))

	s.Equal(http.StatusNotFound, w.Code)
	s.categoryService.AssertExpectations(s.T())
}

func TestAdminCategorySuite(t *testing.T) {
	suite.Run(t, new(AdminCategoryTestSuite))
}
