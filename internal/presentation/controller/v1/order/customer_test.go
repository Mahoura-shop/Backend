package order

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	orderdto "github.com/Mahoura-shop/Backend/internal/application/dto/order"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

func init() { gin.SetMode(gin.TestMode) }

const orderUserID = uint(42)
const testOrderID = uint(10)

type CustomerOrderTestSuite struct {
	suite.Suite
	constants    *bootstrap.Constants
	orderService *mocks.OrderServiceMock
	controller   *CustomerOrderController
}

func (s *CustomerOrderTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.orderService = mocks.NewOrderServiceMock()
	s.controller = NewCustomerOrderController(s.constants, s.orderService)
}

func (s *CustomerOrderTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Set(s.constants.Context.ID, orderUserID)
		c.Next()
	})
	if withRecovery {
		r.Use(orderPanicToHTTP())
	}
	r.POST("/order", s.controller.RegisterOrder)
	r.GET("/order", s.controller.GetMyOrders)
	r.GET("/order/pay/verify", s.controller.VerifyPayment)
	r.GET("/order/:orderID", s.controller.GetMyOrderDetail)
	r.POST("/order/:orderID/pay/wallet", s.controller.PayByWallet)
	r.POST("/order/:orderID/pay/gateway", s.controller.InitiatePayment)
	r.GET("/order/:orderID/instalments", s.controller.GetMyOrderInstalments)
	return r
}

func orderPanicToHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case *exception.AuthError:
					c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
				case exception.ConflictErrors:
					c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
				case exception.NotFoundError:
					c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
				case exception.ValidationErrors:
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "validation"})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

func jsonOrderBody(v any) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

// C-12: POST /order → 200, orderID returned, price locked at checkout time (service responsibility)
func (s *CustomerOrderTestSuite) TestC12_RegisterOrder_Returns200WithOrderID() {
	req := orderdto.CreateOrderRequest{
		UserID:        orderUserID,
		PaymentMethod: enum.PaymentMethodCash,
	}
	s.orderService.On("RegisterOrder", orderUserID, req).Return(uint(10), nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/order",
		jsonOrderBody(map[string]any{"paymentMethod": enum.PaymentMethodCash}))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	s.Equal(float64(10), data["orderID"])
	s.orderService.AssertExpectations(s.T())
}

// C-13: POST /order with empty cart → service returns ConflictError → non-200
func (s *CustomerOrderTestSuite) TestC13_RegisterOrder_EmptyCart_ReturnsError() {
	req := orderdto.CreateOrderRequest{
		UserID:        orderUserID,
		PaymentMethod: enum.PaymentMethodCash,
	}
	s.orderService.On("RegisterOrder", orderUserID, req).Return(uint(0), exception.ConflictErrors{
		Errors: []exception.FieldError{{Field: "cart", Tag: "empty"}},
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/order",
		jsonOrderBody(map[string]any{"paymentMethod": enum.PaymentMethodCash}))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(true).ServeHTTP(w, r)

	s.NotEqual(http.StatusOK, w.Code)
	s.orderService.AssertExpectations(s.T())
}

// C-14: GET /order → 200 with order list
func (s *CustomerOrderTestSuite) TestC14_GetMyOrders_Returns200() {
	s.orderService.On("GetUserOrders", orderUserID).Return([]orderdto.OrderCredential{
		{ID: testOrderID, UserID: orderUserID, Status: enum.OrderStatusPending},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/order", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	orders := body["data"].([]any)
	s.Len(orders, 1)
	s.orderService.AssertExpectations(s.T())
}

// C-15: GET /order/:orderID → 200 with correct items + status
func (s *CustomerOrderTestSuite) TestC15_GetOrderDetail_Returns200() {
	s.orderService.On("GetOrder", testOrderID).Return(&orderdto.OrderCredential{
		ID:     testOrderID,
		UserID: orderUserID,
		Status: enum.OrderStatusPending,
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/order/10", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	s.Equal(float64(testOrderID), data["id"])
	s.orderService.AssertExpectations(s.T())
}

// C-16: POST /order/:orderID/pay/wallet sufficient balance → 200
func (s *CustomerOrderTestSuite) TestC16_PayByWallet_SufficientBalance_Returns200() {
	s.orderService.On("PayOrderByWallet", orderUserID, testOrderID).Return(nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/order/10/pay/wallet", nil))

	s.Equal(http.StatusOK, w.Code)
	s.orderService.AssertExpectations(s.T())
}

// C-17: POST /order/:orderID/pay/wallet insufficient balance → service error → non-200
func (s *CustomerOrderTestSuite) TestC17_PayByWallet_InsufficientBalance_ReturnsError() {
	s.orderService.On("PayOrderByWallet", orderUserID, testOrderID).Return(
		exception.ConflictErrors{Errors: []exception.FieldError{{Field: "wallet", Tag: "insufficientBalance"}}},
	)

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/order/10/pay/wallet", nil))

	s.NotEqual(http.StatusOK, w.Code)
	s.orderService.AssertExpectations(s.T())
}

// C-18: POST /order/:orderID/pay/gateway → 200 with gatewayURL
func (s *CustomerOrderTestSuite) TestC18_InitiateGatewayPayment_Returns200WithURL() {
	s.orderService.On("InitiateGatewayPayment", orderUserID, testOrderID).Return(
		&orderdto.PaymentGatewayResponse{GatewayURL: "https://gateway.example.com/pay/abc"},
		nil,
	)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/order/10/pay/gateway", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	s.NotEmpty(data["gatewayURL"])
	s.orderService.AssertExpectations(s.T())
}

// C-19: GET /order/pay/verify?Authority=X&Status=OK → 200 (gateway success callback)
func (s *CustomerOrderTestSuite) TestC19_VerifyPayment_Success_Returns200() {
	s.orderService.On("VerifyGatewayPayment", "auth-token-abc", "OK").Return(nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/order/pay/verify?Authority=auth-token-abc&Status=OK", nil))

	s.Equal(http.StatusOK, w.Code)
	s.orderService.AssertExpectations(s.T())
}

// C-20: GET /order/pay/verify?Status=NOK → service returns error → non-200
func (s *CustomerOrderTestSuite) TestC20_VerifyPayment_Failed_ReturnsError() {
	s.orderService.On("VerifyGatewayPayment", "auth-token-abc", "NOK").Return(
		exception.ConflictErrors{Errors: []exception.FieldError{{Field: "payment", Tag: "failed"}}},
	)

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/order/pay/verify?Authority=auth-token-abc&Status=NOK", nil))

	s.NotEqual(http.StatusOK, w.Code)
	s.orderService.AssertExpectations(s.T())
}

// C-21: Pay already-paid order → service returns ConflictError → non-200
func (s *CustomerOrderTestSuite) TestC21_PayAlreadyPaidOrder_ReturnsError() {
	s.orderService.On("PayOrderByWallet", orderUserID, testOrderID).Return(
		exception.ConflictErrors{Errors: []exception.FieldError{{Field: "order", Tag: "alreadyPaid"}}},
	)

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/order/10/pay/wallet", nil))

	s.NotEqual(http.StatusOK, w.Code)
	s.orderService.AssertExpectations(s.T())
}

// C-15b: GET /order/:orderID instalments → 200 (C-M2-04 overlap tested here for customer)
func (s *CustomerOrderTestSuite) TestC_GetInstalments_Returns200() {
	s.orderService.On("GetOrderInstalments", testOrderID).Return(
		[]orderdto.InstalmentCredential{{ID: 1, OrderID: testOrderID, Number: 1, Amount: 100_000}},
		nil,
	)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/order/10/instalments", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	instalments := body["data"].([]any)
	s.Len(instalments, 1)
	s.orderService.AssertExpectations(s.T())
}

func TestCustomerOrderSuite(t *testing.T) {
	suite.Run(t, new(CustomerOrderTestSuite))
}
