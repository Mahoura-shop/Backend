package brand

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	branddto "github.com/Mahoura-shop/Backend/internal/application/dto/brand"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

func init() { gin.SetMode(gin.TestMode) }

const brandID = uint(5)

type AdminBrandTestSuite struct {
	suite.Suite
	constants    *bootstrap.Constants
	brandService *mocks.BrandServiceMock
	controller   *AdminBrandController
}

func (s *AdminBrandTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.brandService = mocks.NewBrandServiceMock()
	s.controller = NewAdminBrandController(s.constants, &bootstrap.Pagination{}, s.brandService)
}

func (s *AdminBrandTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Next()
	})
	if withRecovery {
		r.Use(brandPanicToHTTP())
	}
	r.GET("/admin/brands", s.controller.GetBrands)
	r.POST("/admin/brands", s.controller.CreateBrand)
	r.PATCH("/admin/brands/:brandID", s.controller.UpdateBrand)
	r.DELETE("/admin/brands/:brandID", s.controller.DeleteBrand)
	return r
}

func brandPanicToHTTP() gin.HandlerFunc {
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

func brandMultipartBody(fields map[string]string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for k, v := range fields {
		_ = w.WriteField(k, v)
	}
	w.Close()
	return body, w.FormDataContentType()
}

// ADM-04: POST /admin/brands → 200, brand saved
func (s *AdminBrandTestSuite) TestADM04_CreateBrand_Returns200() {
	req := branddto.CreateBrandRequest{
		Name:     "Nivea",
		Slug:     "nivea",
		IsActive: true,
	}
	s.brandService.On("CreateBrand", req).Return(nil)

	body, ct := brandMultipartBody(map[string]string{"name": "Nivea", "slug": "nivea"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/brands", body)
	r.Header.Set("Content-Type", ct)
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.brandService.AssertExpectations(s.T())
}

// ADM-04b: POST /admin/brands missing name → 422
func (s *AdminBrandTestSuite) TestADM04b_CreateBrand_MissingName_Returns422() {
	body, ct := brandMultipartBody(map[string]string{"slug": "nivea"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/brands", body)
	r.Header.Set("Content-Type", ct)
	s.newRouter(true).ServeHTTP(w, r)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
	s.brandService.AssertNotCalled(s.T(), "CreateBrand")
}

// ADM-05: PATCH /admin/brands/:brandID → 200, brand updated
func (s *AdminBrandTestSuite) TestADM05_UpdateBrand_Returns200() {
	newName := "Nivea Pro"
	req := branddto.UpdateBrandRequest{
		ID:   brandID,
		Name: &newName,
	}
	s.brandService.On("UpdateBrand", req).Return(nil)

	body, ct := brandMultipartBody(map[string]string{"name": newName})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/admin/brands/5", body)
	r.Header.Set("Content-Type", ct)
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.brandService.AssertExpectations(s.T())
}

// ADM-06: DELETE /admin/brands/:brandID → 200, brand removed
func (s *AdminBrandTestSuite) TestADM06_DeleteBrand_Returns200() {
	s.brandService.On("DeleteBrand", brandID).Return(nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodDelete, "/admin/brands/5", nil))

	s.Equal(http.StatusOK, w.Code)
	s.brandService.AssertExpectations(s.T())
}

// ADM-06b: Delete non-existent brand → service returns NotFoundError → 404
func (s *AdminBrandTestSuite) TestADM06b_DeleteBrand_NotFound_Returns404() {
	s.brandService.On("DeleteBrand", brandID).Return(exception.NotFoundError{Item: "brand"})

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w,
		httptest.NewRequest(http.MethodDelete, "/admin/brands/5", nil))

	s.Equal(http.StatusNotFound, w.Code)
	s.brandService.AssertExpectations(s.T())
}

func TestAdminBrandSuite(t *testing.T) {
	suite.Run(t, new(AdminBrandTestSuite))
}
