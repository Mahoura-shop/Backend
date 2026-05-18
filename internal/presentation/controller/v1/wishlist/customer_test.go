package wishlist

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	wishlistdto "github.com/Mahoura-shop/Backend/internal/application/dto/wishlist"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

func init() { gin.SetMode(gin.TestMode) }

const wishlistUserID = uint(20)
const wishlistProductID = uint(5)

type CustomerWishlistTestSuite struct {
	suite.Suite
	constants       *bootstrap.Constants
	wishlistService *mocks.WishlistServiceMock
	controller      *CustomerWishlistController
}

func (s *CustomerWishlistTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.wishlistService = mocks.NewWishlistServiceMock()
	s.controller = NewCustomerWishlistController(s.constants, s.wishlistService)
}

func (s *CustomerWishlistTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Set(s.constants.Context.ID, wishlistUserID)
		c.Next()
	})
	if withRecovery {
		r.Use(wishlistPanicToHTTP())
	}
	r.GET("/wishlist", s.controller.GetMyWishlist)
	r.POST("/wishlist/:productID", s.controller.AddToWishlist)
	r.DELETE("/wishlist/:productID", s.controller.RemoveFromWishlist)
	return r
}

func wishlistPanicToHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case exception.ConflictErrors:
					c.JSON(http.StatusConflict, gin.H{"error": e.Error()})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

func (s *CustomerWishlistTestSuite) TestGetMyWishlist_Returns200() {
	s.wishlistService.On("GetMyWishlist", wishlistUserID).Return([]wishlistdto.WishlistItemCredential{}, nil).Once()

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/wishlist", nil))

	s.Equal(http.StatusOK, w.Code)
	s.wishlistService.AssertExpectations(s.T())
}

func (s *CustomerWishlistTestSuite) TestAddToWishlist_Returns201() {
	s.wishlistService.On("AddToWishlist", wishlistUserID, wishlistProductID).Return(nil).Once()

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/wishlist/5", nil))

	s.Equal(http.StatusCreated, w.Code)
	s.wishlistService.AssertExpectations(s.T())
}

func (s *CustomerWishlistTestSuite) TestAddToWishlist_Duplicate_NonOK() {
	conflict := exception.ConflictErrors{}
	conflict.Add("wishlist", "alreadyExist")
	s.wishlistService.On("AddToWishlist", wishlistUserID, wishlistProductID).Return(conflict).Once()

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/wishlist/5", nil))

	s.Equal(http.StatusConflict, w.Code)
	s.wishlistService.AssertExpectations(s.T())
}

func (s *CustomerWishlistTestSuite) TestRemoveFromWishlist_Returns200() {
	s.wishlistService.On("RemoveFromWishlist", wishlistUserID, wishlistProductID).Return(nil).Once()

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/wishlist/5", nil))

	s.Equal(http.StatusOK, w.Code)
	s.wishlistService.AssertExpectations(s.T())
}

func TestCustomerWishlistSuite(t *testing.T) {
	suite.Run(t, new(CustomerWishlistTestSuite))
}
