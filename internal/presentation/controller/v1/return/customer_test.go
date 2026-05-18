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

func init() { gin.SetMode(gin.TestMode) }

const returnUserID = uint(10)
const returnItemID = uint(3)

type CustomerReturnTestSuite struct {
	suite.Suite
	constants     *bootstrap.Constants
	returnService *mocks.ReturnServiceMock
	controller    *CustomerReturnController
}

func (s *CustomerReturnTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.returnService = mocks.NewReturnServiceMock()
	s.controller = NewCustomerReturnController(s.constants, s.returnService)
}

func (s *CustomerReturnTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Set(s.constants.Context.ID, returnUserID)
		c.Next()
	})
	if withRecovery {
		r.Use(returnPanicToHTTP())
	}
	r.POST("/returns", s.controller.RequestReturn)
	r.GET("/returns", s.controller.GetMyReturns)
	return r
}

func returnPanicToHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case exception.NotFoundError:
					c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
				case exception.ForbiddenError:
					c.JSON(http.StatusForbidden, gin.H{"error": e.Error()})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

func (s *CustomerReturnTestSuite) TestGetMyReturns_Returns200() {
	s.returnService.On("GetMyReturns", returnUserID).Return([]returndto.ReturnCredential{}, nil).Once()

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/returns", nil))

	s.Equal(http.StatusOK, w.Code)
	s.returnService.AssertExpectations(s.T())
}

func (s *CustomerReturnTestSuite) TestRequestReturn_Returns201() {
	req := returndto.RequestReturnRequest{
		UserID:      returnUserID,
		OrderItemID: returnItemID,
		Reason:      "damaged",
		Quantity:    1,
	}
	s.returnService.On("RequestReturn", req).Return(nil).Once()

	body, _ := json.Marshal(map[string]any{
		"orderItemID": returnItemID,
		"reason":      "damaged",
		"quantity":    1,
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/returns", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusCreated, w.Code)
	s.returnService.AssertExpectations(s.T())
}

func (s *CustomerReturnTestSuite) TestRequestReturn_ServiceError_NonOK() {
	req := returndto.RequestReturnRequest{
		UserID:      returnUserID,
		OrderItemID: returnItemID,
		Reason:      "damaged",
		Quantity:    1,
	}
	s.returnService.On("RequestReturn", req).Return(exception.NotFoundError{Item: "order item"}).Once()

	body, _ := json.Marshal(map[string]any{
		"orderItemID": returnItemID,
		"reason":      "damaged",
		"quantity":    1,
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/returns", bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(true).ServeHTTP(w, r)

	s.Equal(http.StatusNotFound, w.Code)
	s.returnService.AssertExpectations(s.T())
}

func TestCustomerReturnSuite(t *testing.T) {
	suite.Run(t, new(CustomerReturnTestSuite))
}
