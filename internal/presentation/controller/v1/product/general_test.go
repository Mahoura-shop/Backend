package product

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	productdto "github.com/Mahoura-shop/Backend/internal/application/dto/product"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type GuestProductTestSuite struct {
	suite.Suite
	constants      *bootstrap.Constants
	productService *mocks.ProductServiceMock
	userRepo       *mocks.UserRepositoryMock
	dbMock         *mocks.DatabaseMock
	controller     *GeneralProductController
	router         *gin.Engine
}

func (s *GuestProductTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.productService = mocks.NewProductServiceMock()
	s.userRepo = mocks.NewUserRepositoryMock()
	s.dbMock = mocks.NewDatabaseMock()

	s.controller = NewGeneralProductController(
		s.constants,
		s.productService,
		s.userRepo,
		s.dbMock,
	)

	s.router = gin.New()
	s.router.GET("/products", s.controller.GetProducts)
	s.router.GET("/products/slug/:slug", s.controller.GetProductBySlug)
}

func sampleProduct() productdto.ProductCredential {
	return productdto.ProductCredential{
		ID:            1,
		Name:          "Test Cream",
		Slug:          "test-cream",
		ConsumerPrice: 500_000,
		Step1Price:    200_000,
		Step2Price:    300_000,
		Step3Price:    380_000,
		Step4Price:    450_000,
		CategoryID:    10,
		BrandID:       5,
		Quantity:      100,
		MinOrder:      1,
		IsActive:      true,
	}
}

// G-01: GET /products without auth returns 200 and a product list
func (s *GuestProductTestSuite) TestG01_GetProducts_Returns200() {
	product := sampleProduct()
	s.productService.On("SearchProductsWithPagination", productdto.ProductFilterRequest{}).
		Return(&productdto.ProductSearchResponse{
			Products:   []productdto.ProductCredential{product},
			TotalCount: 1,
		}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	s.Equal(float64(200), body["statusCode"])
	s.productService.AssertExpectations(s.T())
}

// G-01b: Returned products list is non-empty
func (s *GuestProductTestSuite) TestG01b_GetProducts_ReturnsProducts() {
	product := sampleProduct()
	s.productService.On("SearchProductsWithPagination", productdto.ProductFilterRequest{}).
		Return(&productdto.ProductSearchResponse{
			Products:   []productdto.ProductCredential{product},
			TotalCount: 1,
		}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	s.router.ServeHTTP(w, req)

	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))

	data, ok := body["data"].(map[string]any)
	s.Require().True(ok, "data field must be an object")

	products, ok := data["products"].([]any)
	s.Require().True(ok, "data.products must be an array")
	s.Greater(len(products), 0, "product list must not be empty")
}

// G-02: GET /products/slug/:slug returns 200 with consumerPrice in response
func (s *GuestProductTestSuite) TestG02_GetProductBySlug_Returns200WithConsumerPrice() {
	product := sampleProduct()
	s.productService.On("GetProductBySlug", "test-cream").Return(&product, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/products/slug/test-cream", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))

	data, ok := body["data"].(map[string]any)
	s.Require().True(ok, "data must be an object")

	consumerPrice, ok := data["consumerPrice"]
	s.True(ok, "consumerPrice field must be present in product response")
	s.NotNil(consumerPrice, "consumerPrice must not be nil")
	s.productService.AssertExpectations(s.T())
}

// G-02c: Guest resolvedPrice must equal ConsumerPrice (not a lower step price)
// NOTE: This test exposes a bug — GetProducts defaults to UserTypeCustomer (Step4Price)
// for unauthenticated requests instead of UserTypeGuest (ConsumerPrice).
// Expected: resolvedPrice == 500_000 (ConsumerPrice)
// If failing: controller uses Step4Price (450_000) for unauthenticated users.
func (s *GuestProductTestSuite) TestG02c_GuestResolvedPriceIsConsumerPrice() {
	product := sampleProduct()
	s.productService.On("SearchProductsWithPagination", productdto.ProductFilterRequest{}).
		Return(&productdto.ProductSearchResponse{
			Products:   []productdto.ProductCredential{product},
			TotalCount: 1,
		}, nil)

	w := httptest.NewRecorder()
	// No Authorization header → no userID set in context → guest
	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	s.router.ServeHTTP(w, req)

	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))

	data := body["data"].(map[string]any)
	products := data["products"].([]any)
	s.Require().NotEmpty(products)

	firstProduct := products[0].(map[string]any)
	resolvedPrice := firstProduct["resolvedPrice"].(float64)
	consumerPrice := firstProduct["consumerPrice"].(float64)

	s.Equal(consumerPrice, resolvedPrice,
		"guest resolvedPrice must equal consumerPrice (500000), got %v — check GetProducts: default userType should be UserTypeGuest not UserTypeCustomer",
		resolvedPrice)
}

// G-06: GET /products?categoryID=10 passes categoryID filter to service
func (s *GuestProductTestSuite) TestG06_GetProducts_FilterByCategory() {
	categoryID := uint(10)
	product := sampleProduct()

	s.productService.On("SearchProductsWithPagination", productdto.ProductFilterRequest{
		CategoryID: &categoryID,
	}).Return(&productdto.ProductSearchResponse{
		Products:   []productdto.ProductCredential{product},
		TotalCount: 1,
	}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/products?categoryID=10", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.productService.AssertExpectations(s.T())
}

// G-07: GET /products?brandID=5 passes brandID filter to service
func (s *GuestProductTestSuite) TestG07_GetProducts_FilterByBrand() {
	brandID := uint(5)
	product := sampleProduct()

	s.productService.On("SearchProductsWithPagination", productdto.ProductFilterRequest{
		BrandID: &brandID,
	}).Return(&productdto.ProductSearchResponse{
		Products:   []productdto.ProductCredential{product},
		TotalCount: 1,
	}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/products?brandID=5", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.productService.AssertExpectations(s.T())
}

func TestGuestProductSuite(t *testing.T) {
	suite.Run(t, new(GuestProductTestSuite))
}
