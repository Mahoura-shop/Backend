package product

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

const adminProductID = uint(1)
const adminCategoryID = uint(10)

type AdminProductTestSuite struct {
	suite.Suite
	constants      *bootstrap.Constants
	productService *mocks.ProductServiceMock
	controller     *AdminProductController
}

func (s *AdminProductTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.productService = mocks.NewProductServiceMock()
	s.controller = NewAdminProductController(s.constants, &bootstrap.Pagination{}, s.productService)
}

func (s *AdminProductTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Next()
	})
	if withRecovery {
		r.Use(adminProductPanicToHTTP())
	}
	r.GET("/admin/products", s.controller.GetProducts)
	r.GET("/admin/products/slug/:slug", s.controller.GetProduct)
	r.GET("/admin/products/category/:categoryID", s.controller.GetCategoryProducts)
	r.GET("/admin/products/prices", s.controller.GetProductPrices)
	r.POST("/admin/products", s.controller.CreateProduct)
	r.PATCH("/admin/products/:productID", s.controller.UpdateProduct)
	r.DELETE("/admin/products/:productID", s.controller.DeleteProduct)
	r.PUT("/admin/products/prices", s.controller.UpdateProductPrices)
	return r
}

func adminProductPanicToHTTP() gin.HandlerFunc {
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

func productMultipartBody(fields map[string]string) (*bytes.Buffer, string) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	for k, v := range fields {
		_ = w.WriteField(k, v)
	}
	w.Close()
	return body, w.FormDataContentType()
}

// ADM-07: POST /admin/products (all prices set) → 200, product saved
func (s *AdminProductTestSuite) TestADM07_CreateProduct_Returns200() {
	req := productdto.CreateProductRequest{
		Name:  "Cream X",
		Slug:  "cream-x",
		Price: 10.5,
	}
	s.productService.On("CreateProduct", req).Return(nil)

	body, ct := productMultipartBody(map[string]string{
		"name":  "Cream X",
		"slug":  "cream-x",
		"price": "10.5",
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/products", body)
	r.Header.Set("Content-Type", ct)
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.productService.AssertExpectations(s.T())
}

// ADM-07b: POST /admin/products missing required name → 422
func (s *AdminProductTestSuite) TestADM07b_CreateProduct_MissingName_Returns422() {
	body, ct := productMultipartBody(map[string]string{"slug": "cream-x", "price": "10.5"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/products", body)
	r.Header.Set("Content-Type", ct)
	s.newRouter(true).ServeHTTP(w, r)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
	s.productService.AssertNotCalled(s.T(), "CreateProduct")
}

// ADM-09: PATCH /admin/products/:productID → 200, product updated
func (s *AdminProductTestSuite) TestADM09_UpdateProduct_Returns200() {
	newName := "Cream X Pro"
	req := productdto.UpdateProductRequest{
		ID:   adminProductID,
		Name: &newName,
	}
	s.productService.On("UpdateProduct", req).Return(nil)

	body, ct := productMultipartBody(map[string]string{"name": newName})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/admin/products/1", body)
	r.Header.Set("Content-Type", ct)
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.productService.AssertExpectations(s.T())
}

// ADM-10: DELETE /admin/products/:productID → 200, product removed
func (s *AdminProductTestSuite) TestADM10_DeleteProduct_Returns200() {
	s.productService.On("DeleteProduct", adminProductID).Return(nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodDelete, "/admin/products/1", nil))

	s.Equal(http.StatusOK, w.Code)
	s.productService.AssertExpectations(s.T())
}

// ADM-10b: Delete non-existent product → 404
func (s *AdminProductTestSuite) TestADM10b_DeleteProduct_NotFound_Returns404() {
	s.productService.On("DeleteProduct", adminProductID).Return(exception.NotFoundError{Item: "product"})

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w,
		httptest.NewRequest(http.MethodDelete, "/admin/products/1", nil))

	s.Equal(http.StatusNotFound, w.Code)
	s.productService.AssertExpectations(s.T())
}

// ADM-11: GET /admin/products/slug/:slug → 200 with correct product
func (s *AdminProductTestSuite) TestADM11_GetProductBySlug_Returns200() {
	product := sampleProduct()
	s.productService.On("GetProductBySlug", "test-cream").Return(&product, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/admin/products/slug/test-cream", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	s.Equal("test-cream", data["slug"])
	s.productService.AssertExpectations(s.T())
}

// ADM-12: GET /admin/products/category/:categoryID → 200, filtered list
func (s *AdminProductTestSuite) TestADM12_GetCategoryProducts_Returns200() {
	s.productService.On("GetCategoryProducts", adminCategoryID).Return(
		[]productdto.ProductCredential{sampleProduct()}, nil,
	)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/admin/products/category/10", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	products := body["data"].([]any)
	s.Len(products, 1)
	s.productService.AssertExpectations(s.T())
}

// ADM-13: GET /admin/products/prices → 200 with all 4 step prices per product
func (s *AdminProductTestSuite) TestADM13_GetProductPrices_Returns200() {
	s.productService.On("GetProductPrices").Return([]productdto.ProductPrices{
		{ID: 1, Name: "Cream X", IRRPrice: 300_000},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/admin/products/prices", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	prices := body["data"].([]any)
	s.Len(prices, 1)
	s.productService.AssertExpectations(s.T())
}

// ADM-14: PUT /admin/products/prices → 200, prices updated
func (s *AdminProductTestSuite) TestADM14_UpdateProductPrices_Returns200() {
	s.productService.On("UpdateProductsPrice", []productdto.ProductPriceUpdateCredentials{
		{ID: 1, IRRPrice: 350_000},
	}).Return(nil)

	payload, _ := json.Marshal(map[string]any{
		"productPrices": []map[string]any{
			{"id": 1, "irrPrice": 350_000},
		},
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/admin/products/prices", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.productService.AssertExpectations(s.T())
}

func TestAdminProductSuite(t *testing.T) {
	suite.Run(t, new(AdminProductTestSuite))
}
