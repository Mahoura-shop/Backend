package review

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type AdminReviewTestSuite struct {
	suite.Suite
	constants     *bootstrap.Constants
	reviewService *mocks.ReviewServiceMock
	controller    *AdminReviewController
}

func (s *AdminReviewTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.reviewService = mocks.NewReviewServiceMock()
	s.controller = NewAdminReviewController(s.constants, s.reviewService)
}

func (s *AdminReviewTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Next()
	})
	if withRecovery {
		r.Use(reviewPanicToHTTP())
	}
	r.DELETE("/admin/reviews/:reviewID", s.controller.DeleteReview)
	return r
}

func (s *AdminReviewTestSuite) TestDeleteReview_Returns200() {
	s.reviewService.On("DeleteReview", uint(3)).Return(nil).Once()

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/admin/reviews/3", nil))

	s.Equal(http.StatusOK, w.Code)
	s.reviewService.AssertExpectations(s.T())
}

func (s *AdminReviewTestSuite) TestDeleteReview_NotFound_NonOK() {
	s.reviewService.On("DeleteReview", uint(99)).Return(exception.NotFoundError{Item: "review"}).Once()

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w, httptest.NewRequest(http.MethodDelete, "/admin/reviews/99", nil))

	s.NotEqual(http.StatusOK, w.Code)
	s.reviewService.AssertExpectations(s.T())
}

func TestAdminReviewSuite(t *testing.T) {
	suite.Run(t, new(AdminReviewTestSuite))
}
