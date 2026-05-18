package service

import (
	"time"

	"github.com/Mahoura-shop/Backend/bootstrap"
	returndto "github.com/Mahoura-shop/Backend/internal/application/dto/return"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type ReturnService struct {
	constants        *bootstrap.Constants
	returnRepository postgres.ReturnRepository
	orderRepository  postgres.OrderRepository
	walletRepository postgres.WalletRepository
	db               database.Database
}

type ReturnServiceDeps struct {
	Constants        *bootstrap.Constants
	ReturnRepository postgres.ReturnRepository
	OrderRepository  postgres.OrderRepository
	WalletRepository postgres.WalletRepository
	DB               database.Database
}

func NewReturnService(deps ReturnServiceDeps) *ReturnService {
	return &ReturnService{
		constants:        deps.Constants,
		returnRepository: deps.ReturnRepository,
		orderRepository:  deps.OrderRepository,
		walletRepository: deps.WalletRepository,
		db:               deps.DB,
	}
}

func (s *ReturnService) toCredential(r *entity.Return) returndto.ReturnCredential {
	cred := returndto.ReturnCredential{
		ID:           r.ID,
		OrderItemID:  r.OrderItemID,
		UserID:       r.UserID,
		Status:       string(r.Status),
		Reason:       r.Reason,
		Quantity:     r.Quantity,
		RefundAmount: r.RefundAmount,
		RequestedAt:  r.RequestedAt,
		ApprovedAt:   r.ApprovedAt,
		RefundedAt:   r.RefundedAt,
	}
	if r.OrderItem.Product.Name != "" {
		cred.ProductName = r.OrderItem.Product.Name
	}
	if r.User.Phone != "" {
		cred.UserPhone = r.User.Phone
	}
	return cred
}

func (s *ReturnService) RequestReturn(req returndto.RequestReturnRequest) error {
	item, err := s.orderRepository.FindOrderItemByID(s.db, req.OrderItemID)
	if err != nil {
		return err
	}
	if item == nil {
		return exception.NotFoundError{Item: "order item"}
	}
	if item.Order.UserID != req.UserID {
		return exception.ForbiddenError{Message: "order item does not belong to user"}
	}
	if req.Quantity > item.Count {
		return exception.ForbiddenError{Message: "quantity exceeds order item count"}
	}

	ret := entity.Return{
		OrderItemID:  req.OrderItemID,
		UserID:       req.UserID,
		Status:       enum.ReturnStatusRequested,
		Reason:       req.Reason,
		Quantity:     req.Quantity,
		RefundAmount: item.PriceSnapshot * req.Quantity,
	}
	_, err = s.returnRepository.CreateReturn(s.db, ret)
	return err
}

func (s *ReturnService) GetMyReturns(userID uint) ([]returndto.ReturnCredential, error) {
	returns, err := s.returnRepository.GetReturnsByUserID(s.db, userID)
	if err != nil {
		return nil, err
	}
	creds := make([]returndto.ReturnCredential, 0, len(returns))
	for _, r := range returns {
		creds = append(creds, s.toCredential(r))
	}
	return creds, nil
}

func (s *ReturnService) GetAllReturns(status string) ([]returndto.ReturnCredential, error) {
	var returns []*entity.Return
	var err error
	if status != "" {
		returns, err = s.returnRepository.GetReturnsByStatus(s.db, enum.ReturnStatus(status))
	} else {
		returns, err = s.returnRepository.GetAllReturns(s.db)
	}
	if err != nil {
		return nil, err
	}
	creds := make([]returndto.ReturnCredential, 0, len(returns))
	for _, r := range returns {
		creds = append(creds, s.toCredential(r))
	}
	return creds, nil
}

func (s *ReturnService) ReviewReturn(returnID uint, req returndto.ReviewReturnRequest) error {
	ret, err := s.returnRepository.FindReturnByID(s.db, returnID)
	if err != nil {
		return err
	}
	if ret == nil {
		return exception.NotFoundError{Item: "return"}
	}
	now := time.Now()
	if req.Action == "approve" {
		ret.Status = enum.ReturnStatusApproved
		ret.ApprovedAt = &now
		ret.ApprovedByID = &req.AdminID
	} else {
		ret.Status = enum.ReturnStatusRejected
	}
	return s.returnRepository.UpdateReturn(s.db, *ret)
}

func (s *ReturnService) ProcessRefund(returnID uint, adminID uint) error {
	ret, err := s.returnRepository.FindReturnByID(s.db, returnID)
	if err != nil {
		return err
	}
	if ret == nil {
		return exception.NotFoundError{Item: "return"}
	}
	if ret.Status != enum.ReturnStatusApproved {
		return exception.ForbiddenError{Message: "return must be approved before refund"}
	}

	if _, err := s.walletRepository.DepositWallet(s.db, ret.UserID, ret.RefundAmount); err != nil {
		return err
	}

	now := time.Now()
	ret.Status = enum.ReturnStatusRefunded
	ret.RefundedAt = &now
	return s.returnRepository.UpdateReturn(s.db, *ret)
}
