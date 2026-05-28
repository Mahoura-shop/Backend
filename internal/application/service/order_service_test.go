package service

import (
	"errors"
	"testing"

	orderdto "github.com/Mahoura-shop/Backend/internal/application/dto/order"
	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/suite"
)

type OrderServiceTestSuite struct {
	suite.Suite
	constants   *bootstrap.Constants
	orderRepo   *mocks.OrderRepositoryMock
	db          *mocks.DatabaseMock
	service     *OrderService
}

func (s *OrderServiceTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.orderRepo = mocks.NewOrderRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewOrderService(OrderServiceDeps{
		Constants:       s.constants,
		OrderRepository: s.orderRepo,
		DB:              s.db,
	})
}

func sampleOrder() entity.Order {
	addrID := uint(10)
	return entity.Order{
		Model:         dbmodel.Model{ID: 1},
		UserID:        42,
		AddressID:     &addrID,
		PaymentMethod: enum.PaymentMethodWallet,
		Status:        enum.OrderStatusPending,
		TotalAmount:   500_000,
		ShippingCost:  20_000,
	}
}

func (s *OrderServiceTestSuite) TestParseOrder_MapsBasicFields() {
	order := sampleOrder()

	result := s.service.ParseOrder(order)

	s.Equal(uint(1), result.ID)
	s.Equal(uint(42), result.UserID)
	s.Require().NotNil(result.AddressID)
	s.Equal(uint(10), *result.AddressID)
	s.Equal(enum.PaymentMethodWallet, result.PaymentMethod)
	s.Equal(enum.OrderStatusPending, result.Status)
	s.Equal(uint(500_000), result.TotalAmount)
	s.Equal(uint(20_000), result.ShippingCost)
}

func (s *OrderServiceTestSuite) TestParseOrder_UserWithNamePopulated() {
	order := sampleOrder()
	order.User = entity.User{
		FirstName: "Ali",
		LastName:  "Rezaei",
		Phone:     "+989123456789",
	}

	result := s.service.ParseOrder(order)

	s.Require().NotNil(result.User)
	s.Equal("Ali", result.User.FirstName)
	s.Equal("Rezaei", result.User.LastName)
}

func (s *OrderServiceTestSuite) TestParseOrder_NoUser_NilUserField() {
	order := sampleOrder()

	result := s.service.ParseOrder(order)

	s.Nil(result.User)
}

func (s *OrderServiceTestSuite) TestGetOrder_Success() {
	order := sampleOrder()
	s.orderRepo.On("FindOrderByID", s.db, uint(1)).Return(&order, nil).Once()

	result, err := s.service.GetOrder(1)

	s.NoError(err)
	s.Require().NotNil(result)
	s.Equal(uint(1), result.ID)
	s.orderRepo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestGetOrder_NotFound() {
	var nilOrder *entity.Order
	s.orderRepo.On("FindOrderByID", s.db, uint(99)).Return(nilOrder, nil).Once()

	result, err := s.service.GetOrder(99)

	s.Nil(result)
	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
}

func (s *OrderServiceTestSuite) TestGetOrder_RepoError() {
	s.orderRepo.On("FindOrderByID", s.db, uint(1)).Return(nil, errors.New("db error")).Once()

	result, err := s.service.GetOrder(1)

	s.Nil(result)
	s.Error(err)
}

func (s *OrderServiceTestSuite) TestGetOrders_ReturnsList() {
	order := sampleOrder()
	s.orderRepo.On("GetOrders", s.db).Return([]*entity.Order{&order}, nil).Once()

	result, err := s.service.GetOrders()

	s.NoError(err)
	s.Len(result, 1)
	s.Equal(uint(1), result[0].ID)
	s.orderRepo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestGetOrders_Empty() {
	s.orderRepo.On("GetOrders", s.db).Return([]*entity.Order{}, nil).Once()

	result, err := s.service.GetOrders()

	s.NoError(err)
	s.Empty(result)
}

func (s *OrderServiceTestSuite) TestGetOrders_RepoError() {
	s.orderRepo.On("GetOrders", s.db).Return(nil, errors.New("db error")).Once()

	_, err := s.service.GetOrders()

	s.Error(err)
}

func (s *OrderServiceTestSuite) TestGetUserOrders_FiltersByUser() {
	order := sampleOrder()
	s.orderRepo.On("GetOrdersByUserID", s.db, uint(42)).Return([]*entity.Order{&order}, nil).Once()

	result, err := s.service.GetUserOrders(42)

	s.NoError(err)
	s.Len(result, 1)
	s.orderRepo.AssertExpectations(s.T())
}

func (s *OrderServiceTestSuite) TestGetUserOrders_RepoError() {
	s.orderRepo.On("GetOrdersByUserID", s.db, uint(42)).Return(nil, errors.New("db error")).Once()

	_, err := s.service.GetUserOrders(42)

	s.Error(err)
}

func (s *OrderServiceTestSuite) TestParseOrder_ReturnsCorrectType() {
	order := sampleOrder()

	result := s.service.ParseOrder(order)

	_, ok := interface{}(result).(orderdto.OrderCredential)
	s.True(ok)
}

func TestOrderServiceSuite(t *testing.T) {
	suite.Run(t, new(OrderServiceTestSuite))
}
