package service

import (
	"errors"
	"testing"

	returndto "github.com/Mahoura-shop/Backend/internal/application/dto/return"
	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ReturnServiceTestSuite struct {
	suite.Suite
	constants    *bootstrap.Constants
	returnRepo   *mocks.ReturnRepositoryMock
	orderRepo    *mocks.OrderRepositoryMock
	walletRepo   *mocks.WalletRepositoryMock
	db           *mocks.DatabaseMock
	service      *ReturnService
}

func (s *ReturnServiceTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.returnRepo = mocks.NewReturnRepositoryMock()
	s.orderRepo = mocks.NewOrderRepositoryMock()
	s.walletRepo = mocks.NewWalletRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewReturnService(ReturnServiceDeps{
		Constants:        s.constants,
		ReturnRepository: s.returnRepo,
		OrderRepository:  s.orderRepo,
		WalletRepository: s.walletRepo,
		DB:               s.db,
	})
}

func sampleOrderItem(userID uint) *entity.OrderItem {
	return &entity.OrderItem{
		Model:         dbmodel.Model{ID: 10},
		Count:         3,
		PriceSnapshot: 100_000,
		Order:         entity.Order{UserID: userID},
	}
}

func (s *ReturnServiceTestSuite) TestRequestReturn_Success() {
	item := sampleOrderItem(42)
	s.orderRepo.On("FindOrderItemByID", s.db, uint(10)).Return(item, nil).Once()

	expectedReturn := entity.Return{
		OrderItemID:  10,
		UserID:       42,
		Status:       enum.ReturnStatusRequested,
		Reason:       "damaged",
		Quantity:     2,
		RefundAmount: 200_000,
	}
	s.returnRepo.On("CreateReturn", s.db, expectedReturn).Return(&entity.Return{}, nil).Once()

	err := s.service.RequestReturn(returndto.RequestReturnRequest{
		OrderItemID: 10,
		UserID:      42,
		Reason:      "damaged",
		Quantity:    2,
	})

	s.NoError(err)
	s.returnRepo.AssertExpectations(s.T())
	s.orderRepo.AssertExpectations(s.T())
}

func (s *ReturnServiceTestSuite) TestRequestReturn_ItemNotFound() {
	var nilItem *entity.OrderItem
	s.orderRepo.On("FindOrderItemByID", s.db, uint(99)).Return(nilItem, nil).Once()

	err := s.service.RequestReturn(returndto.RequestReturnRequest{
		OrderItemID: 99,
		UserID:      42,
		Reason:      "wrong item",
		Quantity:    1,
	})

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.returnRepo.AssertNotCalled(s.T(), "CreateReturn")
}

func (s *ReturnServiceTestSuite) TestRequestReturn_WrongUser_Forbidden() {
	item := sampleOrderItem(99)
	s.orderRepo.On("FindOrderItemByID", s.db, uint(10)).Return(item, nil).Once()

	err := s.service.RequestReturn(returndto.RequestReturnRequest{
		OrderItemID: 10,
		UserID:      42,
		Reason:      "defective",
		Quantity:    1,
	})

	s.Error(err)
	var forbidden exception.ForbiddenError
	s.True(errors.As(err, &forbidden))
	s.returnRepo.AssertNotCalled(s.T(), "CreateReturn")
}

func (s *ReturnServiceTestSuite) TestRequestReturn_QuantityExceeds_Forbidden() {
	item := sampleOrderItem(42)
	s.orderRepo.On("FindOrderItemByID", s.db, uint(10)).Return(item, nil).Once()

	err := s.service.RequestReturn(returndto.RequestReturnRequest{
		OrderItemID: 10,
		UserID:      42,
		Reason:      "defective",
		Quantity:    10,
	})

	s.Error(err)
	var forbidden exception.ForbiddenError
	s.True(errors.As(err, &forbidden))
	s.returnRepo.AssertNotCalled(s.T(), "CreateReturn")
}

func (s *ReturnServiceTestSuite) TestGetMyReturns_ReturnsList() {
	ret := &entity.Return{
		Model:       dbmodel.Model{ID: 1},
		Status:      enum.ReturnStatusRequested,
		Reason:      "broken",
		Quantity:    1,
		RefundAmount: 100_000,
	}
	s.returnRepo.On("GetReturnsByUserID", s.db, uint(42)).Return([]*entity.Return{ret}, nil).Once()

	result, err := s.service.GetMyReturns(42)

	s.NoError(err)
	s.Len(result, 1)
	s.Equal("requested", result[0].Status)
	s.returnRepo.AssertExpectations(s.T())
}

func (s *ReturnServiceTestSuite) TestGetMyReturns_RepoError() {
	s.returnRepo.On("GetReturnsByUserID", s.db, uint(42)).Return(nil, errors.New("db error")).Once()

	result, err := s.service.GetMyReturns(42)

	s.Nil(result)
	s.Error(err)
}

func (s *ReturnServiceTestSuite) TestGetAllReturns_NoFilter() {
	rets := []*entity.Return{{Status: enum.ReturnStatusRequested}}
	s.returnRepo.On("GetAllReturns", s.db).Return(rets, nil).Once()

	result, err := s.service.GetAllReturns("")

	s.NoError(err)
	s.Len(result, 1)
	s.returnRepo.AssertExpectations(s.T())
}

func (s *ReturnServiceTestSuite) TestGetAllReturns_WithStatusFilter() {
	rets := []*entity.Return{{Status: enum.ReturnStatusApproved}}
	s.returnRepo.On("GetReturnsByStatus", s.db, enum.ReturnStatusApproved).Return(rets, nil).Once()

	result, err := s.service.GetAllReturns("approved")

	s.NoError(err)
	s.Len(result, 1)
	s.returnRepo.AssertExpectations(s.T())
}

func (s *ReturnServiceTestSuite) TestReviewReturn_Approve() {
	ret := &entity.Return{
		Model:  dbmodel.Model{ID: 5},
		Status: enum.ReturnStatusRequested,
	}
	s.returnRepo.On("FindReturnByID", s.db, uint(5)).Return(ret, nil).Once()
	s.returnRepo.On("UpdateReturn", s.db, mock.Anything).Return(nil).Once()

	err := s.service.ReviewReturn(5, returndto.ReviewReturnRequest{Action: "approve", AdminID: 1})

	s.NoError(err)
	s.returnRepo.AssertExpectations(s.T())
}

func (s *ReturnServiceTestSuite) TestReviewReturn_NotFound() {
	var nilReturn *entity.Return
	s.returnRepo.On("FindReturnByID", s.db, uint(99)).Return(nilReturn, nil).Once()

	err := s.service.ReviewReturn(99, returndto.ReviewReturnRequest{Action: "approve"})

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
}

func (s *ReturnServiceTestSuite) TestProcessRefund_NotApproved_Forbidden() {
	ret := &entity.Return{
		Model:  dbmodel.Model{ID: 5},
		Status: enum.ReturnStatusRequested,
	}
	s.returnRepo.On("FindReturnByID", s.db, uint(5)).Return(ret, nil).Once()

	err := s.service.ProcessRefund(5, 1)

	s.Error(err)
	var forbidden exception.ForbiddenError
	s.True(errors.As(err, &forbidden))
	s.walletRepo.AssertNotCalled(s.T(), "DepositWallet")
}

func (s *ReturnServiceTestSuite) TestProcessRefund_Success() {
	ret := &entity.Return{
		Model:        dbmodel.Model{ID: 5},
		Status:       enum.ReturnStatusApproved,
		UserID:       42,
		RefundAmount: 200_000,
	}
	s.returnRepo.On("FindReturnByID", s.db, uint(5)).Return(ret, nil).Once()
	s.walletRepo.On("DepositWallet", s.db, uint(42), uint(200_000)).Return(200_000, nil).Once()
	s.returnRepo.On("UpdateReturn", s.db, mock.Anything).Return(nil).Once()

	err := s.service.ProcessRefund(5, 1)

	s.NoError(err)
	s.returnRepo.AssertExpectations(s.T())
	s.walletRepo.AssertExpectations(s.T())
}

func TestReturnServiceSuite(t *testing.T) {
	suite.Run(t, new(ReturnServiceTestSuite))
}
