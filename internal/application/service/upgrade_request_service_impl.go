package service

import (
	"fmt"

	"github.com/Mahoura-shop/Backend/bootstrap"
	upgraderequestdto "github.com/Mahoura-shop/Backend/internal/application/dto/upgrade_request"
	"github.com/Mahoura-shop/Backend/internal/domain/communication"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type UpgradeRequestService struct {
	constants                *bootstrap.Constants
	upgradeRequestRepository postgres.UpgradeRequestRepository
	userAuditLogRepository   postgres.UserAuditLogRepository
	userRepository           postgres.UserRepository
	smsService               communication.SMSService
	db                       database.Database
}

type UpgradeRequestServiceDeps struct {
	Constants                *bootstrap.Constants
	UpgradeRequestRepository postgres.UpgradeRequestRepository
	UserAuditLogRepository   postgres.UserAuditLogRepository
	UserRepository           postgres.UserRepository
	SMSService               communication.SMSService
	DB                       database.Database
}

func NewUpgradeRequestService(deps UpgradeRequestServiceDeps) *UpgradeRequestService {
	return &UpgradeRequestService{
		constants:                deps.Constants,
		upgradeRequestRepository: deps.UpgradeRequestRepository,
		userAuditLogRepository:   deps.UserAuditLogRepository,
		userRepository:           deps.UserRepository,
		smsService:               deps.SMSService,
		db:                       deps.DB,
	}
}

var validUpgradeTypes = map[enum.UserType]bool{
	enum.UserTypeShopkeeperCash:   true,
	enum.UserTypeShopkeeperCheque: true,
	enum.UserTypeFellow:           true,
}

func (s *UpgradeRequestService) parseRequest(req entity.UpgradeRequest) upgraderequestdto.UpgradeRequestCredential {
	cred := upgraderequestdto.UpgradeRequestCredential{
		ID:            req.ID,
		UserID:        req.UserID,
		RequestedType: req.RequestedType.String(),
		BusinessName:  req.BusinessName,
		TaxID:         req.TaxID,
		Status:        req.Status,
		StatusLabel:   req.Status.String(),
		AdminNote:     req.AdminNote,
		ReviewedByID:  req.ReviewedByID,
		CreatedAt:     req.CreatedAt,
	}
	if req.User.Phone != "" {
		cred.UserPhone = req.User.Phone
	}
	return cred
}

func (s *UpgradeRequestService) SubmitUpgradeRequest(req upgraderequestdto.SubmitUpgradeRequestRequest) error {
	if !validUpgradeTypes[req.RequestedType] {
		return exception.ForbiddenError{Message: "invalid upgrade target: must be shopkeeperCash, shopkeeperCheque, or fellow"}
	}

	user, err := s.userRepository.FindUserByID(s.db, req.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return exception.NotFoundError{Item: s.constants.Field.User}
	}
	if user.Type == req.RequestedType {
		return exception.ForbiddenError{Message: "you are already in the requested group"}
	}

	existing, err := s.upgradeRequestRepository.FindPendingByUserID(s.db, req.UserID)
	if err != nil {
		return err
	}
	if existing != nil {
		return exception.ForbiddenError{Message: "you already have a pending upgrade request"}
	}

	_, err = s.upgradeRequestRepository.CreateUpgradeRequest(s.db, entity.UpgradeRequest{
		UserID:        req.UserID,
		RequestedType: req.RequestedType,
		BusinessName:  req.BusinessName,
		TaxID:         req.TaxID,
		Status:        enum.UpgradeRequestStatusPending,
	})
	return err
}

func (s *UpgradeRequestService) GetMyUpgradeRequests(userID uint) ([]upgraderequestdto.UpgradeRequestCredential, error) {
	reqs, err := s.upgradeRequestRepository.GetUpgradeRequestsByUserID(s.db, userID)
	if err != nil {
		return nil, err
	}
	var result []upgraderequestdto.UpgradeRequestCredential
	for _, r := range reqs {
		result = append(result, s.parseRequest(*r))
	}
	return result, nil
}

func (s *UpgradeRequestService) GetAllUpgradeRequests(statusStr string) ([]upgraderequestdto.UpgradeRequestCredential, error) {
	var reqs []*entity.UpgradeRequest
	var err error

	switch statusStr {
	case "pending":
		reqs, err = s.upgradeRequestRepository.GetUpgradeRequestsByStatus(s.db, enum.UpgradeRequestStatusPending)
	case "approved":
		reqs, err = s.upgradeRequestRepository.GetUpgradeRequestsByStatus(s.db, enum.UpgradeRequestStatusApproved)
	case "rejected":
		reqs, err = s.upgradeRequestRepository.GetUpgradeRequestsByStatus(s.db, enum.UpgradeRequestStatusRejected)
	case "infoRequested":
		reqs, err = s.upgradeRequestRepository.GetUpgradeRequestsByStatus(s.db, enum.UpgradeRequestStatusInfoRequested)
	default:
		reqs, err = s.upgradeRequestRepository.GetUpgradeRequests(s.db)
	}
	if err != nil {
		return nil, err
	}

	var result []upgraderequestdto.UpgradeRequestCredential
	for _, r := range reqs {
		result = append(result, s.parseRequest(*r))
	}
	return result, nil
}

func (s *UpgradeRequestService) GetUpgradeRequest(requestID uint) (*upgraderequestdto.UpgradeRequestCredential, error) {
	req, err := s.upgradeRequestRepository.FindUpgradeRequestByID(s.db, requestID)
	if err != nil {
		return nil, err
	}
	if req == nil {
		return nil, exception.NotFoundError{Item: "upgradeRequest"}
	}
	cred := s.parseRequest(*req)
	return &cred, nil
}

func (s *UpgradeRequestService) ReviewUpgradeRequest(requestID uint, req upgraderequestdto.ReviewUpgradeRequestRequest) error {
	upgradeReq, err := s.upgradeRequestRepository.FindUpgradeRequestByID(s.db, requestID)
	if err != nil {
		return err
	}
	if upgradeReq == nil {
		return exception.NotFoundError{Item: "upgradeRequest"}
	}
	if upgradeReq.Status != enum.UpgradeRequestStatusPending && upgradeReq.Status != enum.UpgradeRequestStatusInfoRequested {
		return exception.ForbiddenError{Message: "request is already resolved"}
	}

	user, err := s.userRepository.FindUserByID(s.db, upgradeReq.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return exception.NotFoundError{Item: s.constants.Field.User}
	}

	return s.db.WithTransaction(func(tx database.Database) error {
		upgradeReq.AdminNote = req.AdminNote
		upgradeReq.ReviewedByID = &req.AdminID

		switch req.Action {
		case "approve":
			upgradeReq.Status = enum.UpgradeRequestStatusApproved
			oldType := user.Type
			user.Type = upgradeReq.RequestedType
			if err := s.userRepository.UpdateUser(tx, *user); err != nil {
				return err
			}
			if err := s.userAuditLogRepository.CreateUserAuditLog(tx, entity.UserAuditLog{
				UserID:      user.ID,
				ChangedByID: req.AdminID,
				OldType:     oldType,
				NewType:     upgradeReq.RequestedType,
				Reason:      fmt.Sprintf("upgrade request #%d approved", requestID),
			}); err != nil {
				return err
			}
			go s.smsService.SendMessage(user.Phone, fmt.Sprintf("درخواست ارتقاء حساب شما تایید شد. نوع حساب جدید: %s", upgradeReq.RequestedType.String()))

		case "reject":
			upgradeReq.Status = enum.UpgradeRequestStatusRejected
			go s.smsService.SendMessage(user.Phone, "درخواست ارتقاء حساب شما رد شد.")

		case "info":
			upgradeReq.Status = enum.UpgradeRequestStatusInfoRequested
			go s.smsService.SendMessage(user.Phone, "برای بررسی درخواست ارتقاء حساب شما، به اطلاعات بیشتری نیاز است.")
		}

		return s.upgradeRequestRepository.UpdateUpgradeRequest(tx, *upgradeReq)
	})
}

func (s *UpgradeRequestService) ChangeUserType(userID uint, req upgraderequestdto.ChangeUserTypeRequest) error {
	user, err := s.userRepository.FindUserByID(s.db, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return exception.NotFoundError{Item: s.constants.Field.User}
	}

	newType := enum.UserType(req.NewType)
	oldType := user.Type

	return s.db.WithTransaction(func(tx database.Database) error {
		user.Type = newType
		if err := s.userRepository.UpdateUser(tx, *user); err != nil {
			return err
		}
		return s.userAuditLogRepository.CreateUserAuditLog(tx, entity.UserAuditLog{
			UserID:      user.ID,
			ChangedByID: req.AdminID,
			OldType:     oldType,
			NewType:     newType,
			Reason:      req.Reason,
		})
	})
}

func (s *UpgradeRequestService) GetUserAuditLogs(userID uint) ([]upgraderequestdto.UserAuditLogCredential, error) {
	logs, err := s.userAuditLogRepository.GetAuditLogsByUserID(s.db, userID)
	if err != nil {
		return nil, err
	}
	var result []upgraderequestdto.UserAuditLogCredential
	for _, l := range logs {
		result = append(result, upgraderequestdto.UserAuditLogCredential{
			ID:          l.ID,
			UserID:      l.UserID,
			ChangedByID: l.ChangedByID,
			OldType:     l.OldType.String(),
			NewType:     l.NewType.String(),
			Reason:      l.Reason,
			CreatedAt:   l.CreatedAt,
		})
	}
	return result, nil
}
