package usecase

import (
	orderdto "github.com/Mahoura-shop/Backend/internal/application/dto/order"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
)

type OrderService interface {
	ParseOrder(entity.Order) orderdto.OrderCredential
	GetOrder(uint) (*orderdto.OrderCredential, error)
	GetOrders() ([]orderdto.OrderCredential, error)
	GetUserOrders(userID uint) ([]orderdto.OrderCredential, error)
	RegisterOrder(userID uint, req orderdto.CreateOrderRequest) (uint, error)
	UpdateOrderStatus(orderID uint, req orderdto.UpdateOrderStatusRequest) error
	PayOrderByWallet(userID, orderID uint) error
	InitiateGatewayPayment(userID, orderID uint) (*orderdto.PaymentGatewayResponse, error)
	VerifyGatewayPayment(authority string, status string) error
	CancelOrder(orderID uint, reason string) error
	FlagOrderRefund(orderID uint) error
	GetOrderInstalments(orderID uint) ([]orderdto.InstalmentCredential, error)
	GetOrdersByStatus(status enum.OrderStatus) ([]orderdto.OrderCredential, error)
}
