package review

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	reviewdto "github.com/Mahoura-shop/Backend/internal/application/dto/review"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

const reviewUserID = uint(42)

type CustomerReviewTestSuite struct {
	suite.Suite
	constants     *bootstrap.Constants
	reviewService *mocks.ReviewServiceMock
	controller    *CustomerReviewController
}

func (s *CustomerReviewTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.reviewService = mocks.NewReviewServiceMock()
	s.controller = NewCustomerReviewController(s.constants, s.reviewService)
}

func (s *CustomerReviewTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Set(s.constants.Context.ID, reviewUserID)
		c.Next()
	})
	if withRecovery {
		r.Use(reviewPanicToHTTP())
	}
	r.POST("/products/:productID/reviews", s.controller.SubmitReview)
	r.GET("/products/:productID/reviews", s.controller.GetProductReviews)
	return r
}

func reviewPanicToHTTP() gin.HandlerFunc {
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

func (s *CustomerReviewTestSuite) TestSubmitReview_Returns201() {
	req := reviewdto.SubmitReviewRequest{
		UserID:    reviewUserID,
		ProductID: 10,
		Rating:    5,
		Comment:   "Excellent",
	}
	s.reviewService.On("SubmitReview", req).Return(nil).Once()

	body, _ := json.Marshal(map[string]any{"rating": 5, "comment": "Excellent"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/products/10/reviews", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusCreated, w.Code)
	s.reviewService.AssertExpectations(s.T())
}

func (s *CustomerReviewTestSuite) TestSubmitReview_Duplicate_Conflict() {
	req := reviewdto.SubmitReviewRequest{
		UserID:    reviewUserID,
		ProductID: 10,
		Rating:    5,
		Comment:   "Excellent",
	}
	conflict := exception.ConflictErrors{}
	conflict.Add("review", "alreadyExist")
	s.reviewService.On("SubmitReview", req).Return(conflict).Once()

	body, _ := json.Marshal(map[string]any{"rating": 5, "comment": "Excellent"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/products/10/reviews", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(true).ServeHTTP(w, r)

	s.Equal(http.StatusConflict, w.Code)
	s.reviewService.AssertExpectations(s.T())
}

func (s *CustomerReviewTestSuite) TestGetProductReviews_Returns200() {
	s.reviewService.On("GetProductReviews", uint(10)).Return([]reviewdto.ReviewCredential{}, nil).Once()

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/products/10/reviews", nil))

	s.Equal(http.StatusOK, w.Code)
	s.reviewService.AssertExpectations(s.T())
}

func TestCustomerReviewSuite(t *testing.T) {
	suite.Run(t, new(CustomerReviewTestSuite))
}
