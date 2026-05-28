package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	postgresrepo "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type OrderRepositoryMock struct {
	mock.Mock
}

func NewOrderRepositoryMock() *OrderRepositoryMock {
	return &OrderRepositoryMock{}
}

func (m *OrderRepositoryMock) FindOrderByID(db database.Database, orderID uint) (*entity.Order, error) {
	args := m.Called(db, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) GetOrders(db database.Database) ([]*entity.Order, error) {
	args := m.Called(db)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) GetOrdersByUserID(db database.Database, userID uint) ([]*entity.Order, error) {
	args := m.Called(db, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) CreateOrder(db database.Database, order entity.Order) (*entity.Order, error) {
	args := m.Called(db, order)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Order), args.Error(1)
}

func (m *OrderRepositoryMock) UpdateOrder(db database.Database, order entity.Order) error {
	args := m.Called(db, order)
	return args.Error(0)
}

func (m *OrderRepositoryMock) DeleteOrderByID(db database.Database, orderID uint) error {
	args := m.Called(db, orderID)
	return args.Error(0)
}

func (m *OrderRepositoryMock) CreateOrderItem(db database.Database, orderItem entity.OrderItem) error {
	args := m.Called(db, orderItem)
	return args.Error(0)
}

func (m *OrderRepositoryMock) CreateOrderStatusHistory(db database.Database, history entity.OrderStatusHistory) error {
	args := m.Called(db, history)
	return args.Error(0)
}

func (m *OrderRepositoryMock) GetRevenuePerTier(db database.Database) ([]postgresrepo.RevenueByTier, error) {
	args := m.Called(db)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgresrepo.RevenueByTier), args.Error(1)
}

func (m *OrderRepositoryMock) GetOrdersPerDay(db database.Database, days int) ([]postgresrepo.OrderByDay, error) {
	args := m.Called(db, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgresrepo.OrderByDay), args.Error(1)
}

func (m *OrderRepositoryMock) GetRevenuePerDay(db database.Database, days int) ([]postgresrepo.RevenueByDay, error) {
	args := m.Called(db, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgresrepo.RevenueByDay), args.Error(1)
}

func (m *OrderRepositoryMock) FindOrderItemByID(db database.Database, orderItemID uint) (*entity.OrderItem, error) {
	args := m.Called(db, orderItemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.OrderItem), args.Error(1)
}

func (m *OrderRepositoryMock) HasOrderItemsForProduct(db database.Database, productID uint) (bool, error) {
	args := m.Called(db, productID)
	return args.Bool(0), args.Error(1)
}

func (m *OrderRepositoryMock) GetProductOrdersPerDay(db database.Database, productID uint, days int) ([]postgresrepo.OrderByDay, error) {
	args := m.Called(db, productID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgresrepo.OrderByDay), args.Error(1)
}

func (m *OrderRepositoryMock) GetCategoryOrdersPerDay(db database.Database, categoryID uint, days int) ([]postgresrepo.OrderByDay, error) {
	args := m.Called(db, categoryID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgresrepo.OrderByDay), args.Error(1)
}

func (m *OrderRepositoryMock) GetBrandOrdersPerDay(db database.Database, brandID uint, days int) ([]postgresrepo.OrderByDay, error) {
	args := m.Called(db, brandID, days)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgresrepo.OrderByDay), args.Error(1)
}

func (m *OrderRepositoryMock) GetProvinceStats(db database.Database) ([]postgresrepo.ProvinceStatRow, error) {
	args := m.Called(db)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]postgresrepo.ProvinceStatRow), args.Error(1)
}
