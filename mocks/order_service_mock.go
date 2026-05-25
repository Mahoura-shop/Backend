package mocks

import (
	orderdto "github.com/Mahoura-shop/Backend/internal/application/dto/order"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/stretchr/testify/mock"
)

type OrderServiceMock struct {
	mock.Mock
}

func NewOrderServiceMock() *OrderServiceMock {
	return &OrderServiceMock{}
}

func (m *OrderServiceMock) ParseOrder(o entity.Order) orderdto.OrderCredential {
	args := m.Called(o)
	return args.Get(0).(orderdto.OrderCredential)
}

func (m *OrderServiceMock) GetOrder(id uint) (*orderdto.OrderCredential, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orderdto.OrderCredential), args.Error(1)
}

func (m *OrderServiceMock) GetOrders() ([]orderdto.OrderCredential, error) {
	args := m.Called()
	return args.Get(0).([]orderdto.OrderCredential), args.Error(1)
}

func (m *OrderServiceMock) GetUserOrders(userID uint) ([]orderdto.OrderCredential, error) {
	args := m.Called(userID)
	return args.Get(0).([]orderdto.OrderCredential), args.Error(1)
}

func (m *OrderServiceMock) RegisterOrder(userID uint, req orderdto.CreateOrderRequest) (uint, error) {
	args := m.Called(userID, req)
	return args.Get(0).(uint), args.Error(1)
}

func (m *OrderServiceMock) UpdateOrderStatus(orderID uint, req orderdto.UpdateOrderStatusRequest) error {
	args := m.Called(orderID, req)
	return args.Error(0)
}

func (m *OrderServiceMock) PayOrderByWallet(userID, orderID uint) error {
	args := m.Called(userID, orderID)
	return args.Error(0)
}

func (m *OrderServiceMock) InitiateGatewayPayment(userID, orderID uint) (*orderdto.PaymentGatewayResponse, error) {
	args := m.Called(userID, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orderdto.PaymentGatewayResponse), args.Error(1)
}

func (m *OrderServiceMock) VerifyGatewayPayment(authority string, status string) error {
	args := m.Called(authority, status)
	return args.Error(0)
}

func (m *OrderServiceMock) CancelOrder(orderID uint, reason string) error {
	args := m.Called(orderID, reason)
	return args.Error(0)
}

func (m *OrderServiceMock) FlagOrderRefund(orderID uint) error {
	args := m.Called(orderID)
	return args.Error(0)
}

func (m *OrderServiceMock) GetOrderInstalments(orderID uint) ([]orderdto.InstalmentCredential, error) {
	args := m.Called(orderID)
	return args.Get(0).([]orderdto.InstalmentCredential), args.Error(1)
}

func (m *OrderServiceMock) GetOrdersByStatus(status enum.OrderStatus) ([]orderdto.OrderCredential, error) {
	args := m.Called(status)
	return args.Get(0).([]orderdto.OrderCredential), args.Error(1)
}
