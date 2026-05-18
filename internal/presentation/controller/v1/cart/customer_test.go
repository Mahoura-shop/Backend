package cart

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	cartdto "github.com/Mahoura-shop/Backend/internal/application/dto/cart"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

func init() { gin.SetMode(gin.TestMode) }

const cartUserID = uint(42)
const testProductID = uint(7)

type CustomerCartTestSuite struct {
	suite.Suite
	constants   *bootstrap.Constants
	cartService *mocks.CartServiceMock
	controller  *CustomerCartController
}

func (s *CustomerCartTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.cartService = mocks.NewCartServiceMock()
	s.controller = NewCustomerCartController(s.constants, s.cartService)
}

func (s *CustomerCartTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Set(s.constants.Context.ID, cartUserID)
		c.Next()
	})
	if withRecovery {
		r.Use(cartPanicToHTTP())
	}
	r.GET("/cart", s.controller.GetUserCart)
	r.POST("/cart/:productID/add", s.controller.AddProductToCart)
	r.POST("/cart/:productID/remove", s.controller.RemoveProductFromCart)
	return r
}

func cartPanicToHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case *exception.AuthError:
					c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
				case exception.ValidationErrors:
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "validation"})
				case exception.ConflictErrors:
					c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
				case exception.NotFoundError:
					c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

// C-11: GET /cart → 200 with items
func (s *CustomerCartTestSuite) TestC11_GetCart_Returns200() {
	s.cartService.On("GetUserCart", cartUserID).Return(cartdto.CartCredential{
		ID:    1,
		Items: []cartdto.CartItemCredential{},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/cart", nil))

	s.Equal(http.StatusOK, w.Code)
	s.cartService.AssertExpectations(s.T())
}

// C-06: POST /cart/:productID/add → 200
func (s *CustomerCartTestSuite) TestC06_AddToCart_Returns200() {
	s.cartService.On("AddProductToCart", cartdto.UpdateProductCountInCart{
		UserID:    cartUserID,
		ProductID: testProductID,
	}).Return(nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/cart/7/add", nil))

	s.Equal(http.StatusOK, w.Code)
	s.cartService.AssertExpectations(s.T())
}

// C-07: Add same product again → service called again, 200 (service handles qty increment)
func (s *CustomerCartTestSuite) TestC07_AddSameProduct_CallsServiceAgain() {
	req := cartdto.UpdateProductCountInCart{UserID: cartUserID, ProductID: testProductID}
	s.cartService.On("AddProductToCart", req).Return(nil).Times(2)

	router := s.newRouter(false)
	for i := 0; i < 2; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/cart/7/add", nil))
		s.Equal(http.StatusOK, w.Code)
	}
	s.cartService.AssertExpectations(s.T())
}

// C-08: Add below MinOrder → service returns ConflictError → non-200
func (s *CustomerCartTestSuite) TestC08_AddToCart_BelowMinOrder_ReturnsError() {
	s.cartService.On("AddProductToCart", cartdto.UpdateProductCountInCart{
		UserID:    cartUserID,
		ProductID: testProductID,
	}).Return(exception.ConflictErrors{
		Errors: []exception.FieldError{{Field: "cartItem", Tag: "belowMinOrder"}},
	})

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/cart/7/add", nil))

	s.NotEqual(http.StatusOK, w.Code)
	s.cartService.AssertExpectations(s.T())
}

// C-09: Add more than stock → service returns ConflictError → non-200
func (s *CustomerCartTestSuite) TestC09_AddToCart_ExceedsStock_ReturnsError() {
	s.cartService.On("AddProductToCart", cartdto.UpdateProductCountInCart{
		UserID:    cartUserID,
		ProductID: testProductID,
	}).Return(exception.ConflictErrors{
		Errors: []exception.FieldError{{Field: "cartItem", Tag: "exceedsStock"}},
	})

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/cart/7/add", nil))

	s.NotEqual(http.StatusOK, w.Code)
	s.cartService.AssertExpectations(s.T())
}

// C-10: POST /cart/:productID/remove → 200
func (s *CustomerCartTestSuite) TestC10_RemoveFromCart_Returns200() {
	s.cartService.On("RemoveProductFromCart", cartdto.UpdateProductCountInCart{
		UserID:    cartUserID,
		ProductID: testProductID,
	}).Return(nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/cart/7/remove", nil))

	s.Equal(http.StatusOK, w.Code)
	s.cartService.AssertExpectations(s.T())
}

func TestCustomerCartSuite(t *testing.T) {
	suite.Run(t, new(CustomerCartTestSuite))
}
