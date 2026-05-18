package returnctrl

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	returndto "github.com/Mahoura-shop/Backend/internal/application/dto/return"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

const adminUserID = uint(1)

type AdminReturnTestSuite struct {
	suite.Suite
	constants     *bootstrap.Constants
	returnService *mocks.ReturnServiceMock
	controller    *AdminReturnController
}

func (s *AdminReturnTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.returnService = mocks.NewReturnServiceMock()
	s.controller = NewAdminReturnController(s.constants, s.returnService)
}

func (s *AdminReturnTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Set(s.constants.Context.ID, adminUserID)
		c.Next()
	})
	if withRecovery {
		r.Use(returnPanicToHTTP())
	}
	r.GET("/admin/returns", s.controller.GetReturns)
	r.POST("/admin/returns/:returnID/review", s.controller.ReviewReturn)
	r.POST("/admin/returns/:returnID/refund", s.controller.ProcessRefund)
	return r
}

func (s *AdminReturnTestSuite) TestGetReturns_Returns200() {
	s.returnService.On("GetAllReturns", "").Return([]returndto.ReturnCredential{}, nil).Once()

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin/returns", nil))

	s.Equal(http.StatusOK, w.Code)
	s.returnService.AssertExpectations(s.T())
}

func (s *AdminReturnTestSuite) TestGetReturns_WithStatusFilter() {
	s.returnService.On("GetAllReturns", "requested").Return([]returndto.ReturnCredential{}, nil).Once()

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin/returns?status=requested", nil))

	s.Equal(http.StatusOK, w.Code)
	s.returnService.AssertExpectations(s.T())
}

func (s *AdminReturnTestSuite) TestReviewReturn_Approve_Returns200() {
	req := returndto.ReviewReturnRequest{AdminID: adminUserID, Action: "approve"}
	s.returnService.On("ReviewReturn", uint(5), req).Return(nil).Once()

	body, _ := json.Marshal(map[string]any{"action": "approve"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/returns/5/review", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.returnService.AssertExpectations(s.T())
}

func (s *AdminReturnTestSuite) TestProcessRefund_NotFound_NonOK() {
	s.returnService.On("ProcessRefund", uint(99), adminUserID).Return(exception.NotFoundError{Item: "return"}).Once()

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/admin/returns/99/refund", nil))

	s.NotEqual(http.StatusOK, w.Code)
	s.returnService.AssertExpectations(s.T())
}

func TestAdminReturnSuite(t *testing.T) {
	suite.Run(t, new(AdminReturnTestSuite))
}
