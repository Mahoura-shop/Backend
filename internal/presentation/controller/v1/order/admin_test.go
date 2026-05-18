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

const adminOrderID = uint(10)
const adminActorID = uint(1)

type AdminOrderTestSuite struct {
	suite.Suite
	constants    *bootstrap.Constants
	orderService *mocks.OrderServiceMock
	controller   *AdminOrderController
}

func (s *AdminOrderTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.orderService = mocks.NewOrderServiceMock()
	s.controller = NewAdminOrderController(s.constants, s.orderService)
}

func (s *AdminOrderTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Set(s.constants.Context.ID, adminActorID)
		c.Next()
	})
	if withRecovery {
		r.Use(adminOrderPanicToHTTP())
	}
	r.GET("/admin/orders", s.controller.GetOrders)
	r.GET("/admin/orders/:orderID", s.controller.GetOrderDetail)
	r.PATCH("/admin/orders/:orderID/status", s.controller.UpdateOrderStatus)
	r.POST("/admin/orders/:orderID/cancel", s.controller.CancelOrder)
	r.POST("/admin/orders/:orderID/refund", s.controller.FlagRefund)
	r.GET("/admin/orders/:orderID/instalments", s.controller.GetOrderInstalments)
	return r
}

func adminOrderPanicToHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case exception.NotFoundError:
					c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
				case exception.ConflictErrors:
					c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
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

// ADM-30: GET /admin/orders → 200 with all orders
func (s *AdminOrderTestSuite) TestADM30_GetOrders_Returns200() {
	s.orderService.On("GetOrders").Return([]orderdto.OrderCredential{
		{ID: adminOrderID, Status: enum.OrderStatusPending},
		{ID: 11, Status: enum.OrderStatusPaid},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/admin/orders", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	orders := body["data"].([]any)
	s.Len(orders, 2)
	s.orderService.AssertExpectations(s.T())
}

// ADM-31: GET /admin/orders?status=pending → 200 with filtered list
func (s *AdminOrderTestSuite) TestADM31_GetOrdersByStatus_Returns200() {
	s.orderService.On("GetOrdersByStatus", enum.OrderStatusPending).Return([]orderdto.OrderCredential{
		{ID: adminOrderID, Status: enum.OrderStatusPending},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/admin/orders?status=pending", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	orders := body["data"].([]any)
	s.Len(orders, 1)
	s.orderService.AssertExpectations(s.T())
}

// ADM-31b: GET /admin/orders?status=paid → 200, different status filter
func (s *AdminOrderTestSuite) TestADM31b_GetOrdersByStatus_Paid_Returns200() {
	s.orderService.On("GetOrdersByStatus", enum.OrderStatusPaid).Return([]orderdto.OrderCredential{
		{ID: 11, Status: enum.OrderStatusPaid},
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/admin/orders?status=paid", nil))

	s.Equal(http.StatusOK, w.Code)
	s.orderService.AssertExpectations(s.T())
}

// ADM-32: GET /admin/orders/:orderID → 200 with order detail
func (s *AdminOrderTestSuite) TestADM32_GetOrderDetail_Returns200() {
	s.orderService.On("GetOrder", adminOrderID).Return(&orderdto.OrderCredential{
		ID:     adminOrderID,
		UserID: orderUserID,
		Status: enum.OrderStatusPending,
	}, nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/admin/orders/10", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	s.Equal(float64(adminOrderID), data["id"])
	s.orderService.AssertExpectations(s.T())
}

// ADM-33: PATCH /admin/orders/:orderID/status → shipped → 200
func (s *AdminOrderTestSuite) TestADM33_UpdateOrderStatus_Shipped_Returns200() {
	req := orderdto.UpdateOrderStatusRequest{
		Status:      enum.OrderStatusShipped,
		Note:        "dispatched",
		ChangedByID: adminActorID,
	}
	s.orderService.On("UpdateOrderStatus", adminOrderID, req).Return(nil)

	payload, _ := json.Marshal(map[string]any{"status": enum.OrderStatusShipped, "note": "dispatched"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/admin/orders/10/status", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.orderService.AssertExpectations(s.T())
}

// ADM-34: PATCH /admin/orders/:orderID/status → delivered → 200
func (s *AdminOrderTestSuite) TestADM34_UpdateOrderStatus_Delivered_Returns200() {
	req := orderdto.UpdateOrderStatusRequest{
		Status:      enum.OrderStatusDelivered,
		ChangedByID: adminActorID,
	}
	s.orderService.On("UpdateOrderStatus", adminOrderID, req).Return(nil)

	payload, _ := json.Marshal(map[string]any{"status": enum.OrderStatusDelivered})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/admin/orders/10/status", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.orderService.AssertExpectations(s.T())
}

// ADM-35: POST /admin/orders/:orderID/cancel → 200, status = cancelled
func (s *AdminOrderTestSuite) TestADM35_CancelOrder_Returns200() {
	s.orderService.On("CancelOrder", adminOrderID).Return(nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/admin/orders/10/cancel", nil))

	s.Equal(http.StatusOK, w.Code)
	s.orderService.AssertExpectations(s.T())
}

// ADM-35b: Cancel non-existent order → 404
func (s *AdminOrderTestSuite) TestADM35b_CancelOrder_NotFound_Returns404() {
	s.orderService.On("CancelOrder", adminOrderID).Return(exception.NotFoundError{Item: "order"})

	w := httptest.NewRecorder()
	s.newRouter(true).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/admin/orders/10/cancel", nil))

	s.Equal(http.StatusNotFound, w.Code)
	s.orderService.AssertExpectations(s.T())
}

// ADM-36: POST /admin/orders/:orderID/refund → 200, order flagged
func (s *AdminOrderTestSuite) TestADM36_FlagRefund_Returns200() {
	s.orderService.On("FlagOrderRefund", adminOrderID).Return(nil)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodPost, "/admin/orders/10/refund", nil))

	s.Equal(http.StatusOK, w.Code)
	s.orderService.AssertExpectations(s.T())
}

// ADM-37: GET /admin/orders/:orderID/instalments → 200 with instalment list
func (s *AdminOrderTestSuite) TestADM37_GetOrderInstalments_Returns200() {
	s.orderService.On("GetOrderInstalments", adminOrderID).Return(
		[]orderdto.InstalmentCredential{
			{ID: 1, OrderID: adminOrderID, Number: 1, Amount: 200_000},
			{ID: 2, OrderID: adminOrderID, Number: 2, Amount: 200_000},
		}, nil,
	)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/admin/orders/10/instalments", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	instalments := body["data"].([]any)
	s.Len(instalments, 2)
	s.orderService.AssertExpectations(s.T())
}

func TestAdminOrderSuite(t *testing.T) {
	suite.Run(t, new(AdminOrderTestSuite))
}
