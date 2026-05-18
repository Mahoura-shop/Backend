package review

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	reviewdto "github.com/Mahoura-shop/Backend/internal/application/dto/review"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

func init() { gin.SetMode(gin.TestMode) }

type GeneralReviewTestSuite struct {
	suite.Suite
	constants     *bootstrap.Constants
	reviewService *mocks.ReviewServiceMock
	controller    *GeneralReviewController
}

func (s *GeneralReviewTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.reviewService = mocks.NewReviewServiceMock()
	s.controller = NewGeneralReviewController(s.constants, s.reviewService)
}

func (s *GeneralReviewTestSuite) newRouter() *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Next()
	})
	r.GET("/products/:productID/reviews", s.controller.GetProductReviews)
	return r
}

func (s *GeneralReviewTestSuite) TestGetProductReviews_Returns200() {
	s.reviewService.On("GetProductReviews", uint(10)).Return([]reviewdto.ReviewCredential{}, nil).Once()

	w := httptest.NewRecorder()
	s.newRouter().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/products/10/reviews", nil))

	s.Equal(http.StatusOK, w.Code)
	s.reviewService.AssertExpectations(s.T())
}

func (s *GeneralReviewTestSuite) TestGetProductReviews_WithReviews_Returns200() {
	reviews := []reviewdto.ReviewCredential{
		{ID: 1, Rating: 5, Comment: "Perfect"},
	}
	s.reviewService.On("GetProductReviews", uint(3)).Return(reviews, nil).Once()

	w := httptest.NewRecorder()
	s.newRouter().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/products/3/reviews", nil))

	s.Equal(http.StatusOK, w.Code)
	s.reviewService.AssertExpectations(s.T())
}

func TestGeneralReviewSuite(t *testing.T) {
	suite.Run(t, new(GeneralReviewTestSuite))
}
