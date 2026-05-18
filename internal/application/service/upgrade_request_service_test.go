package service

import (
	"errors"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	upgraderequestdto "github.com/Mahoura-shop/Backend/internal/application/dto/upgrade_request"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/suite"
)

type UpgradeRequestServiceTestSuite struct {
	suite.Suite
	constants                *bootstrap.Constants
	upgradeRequestRepository *mocks.UpgradeRequestRepositoryMock
	userRepository           *mocks.UserRepositoryMock
	smsService               *mocks.SMSServiceMock
	db                       *mocks.DatabaseMock
	service                  *UpgradeRequestService
}

func (s *UpgradeRequestServiceTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.upgradeRequestRepository = mocks.NewUpgradeRequestRepositoryMock()
	s.userRepository = mocks.NewUserRepositoryMock()
	s.smsService = mocks.NewSMSServiceMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewUpgradeRequestService(UpgradeRequestServiceDeps{
		Constants:                s.constants,
		UpgradeRequestRepository: s.upgradeRequestRepository,
		UserAuditLogRepository:   mocks.NewUserAuditLogRepositoryMock(),
		UserRepository:           s.userRepository,
		SMSService:               s.smsService,
		DB:                       s.db,
	})
}

func (s *UpgradeRequestServiceTestSuite) TestSubmitUpgradeRequest_InvalidType() {
	req := upgraderequestdto.SubmitUpgradeRequestRequest{
		UserID:        1,
		RequestedType: enum.UserTypeAdmin,
	}

	err := s.service.SubmitUpgradeRequest(req)

	s.Error(err)
	var forbidden exception.ForbiddenError
	s.True(errors.As(err, &forbidden))
}

func (s *UpgradeRequestServiceTestSuite) TestSubmitUpgradeRequest_UserNotFound() {
	req := upgraderequestdto.SubmitUpgradeRequestRequest{
		UserID:        99,
		RequestedType: enum.UserTypeShopkeeperCash,
	}
	s.userRepository.On("FindUserByID", s.db, uint(99)).Return((*entity.User)(nil), nil).Once()

	err := s.service.SubmitUpgradeRequest(req)

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.userRepository.AssertExpectations(s.T())
}

func (s *UpgradeRequestServiceTestSuite) TestSubmitUpgradeRequest_AlreadySameType() {
	user := &entity.User{Model: dbmodel.Model{ID: 1}, Type: enum.UserTypeShopkeeperCash}
	req := upgraderequestdto.SubmitUpgradeRequestRequest{
		UserID:        1,
		RequestedType: enum.UserTypeShopkeeperCash,
	}
	s.userRepository.On("FindUserByID", s.db, uint(1)).Return(user, nil).Once()

	err := s.service.SubmitUpgradeRequest(req)

	s.Error(err)
	var forbidden exception.ForbiddenError
	s.True(errors.As(err, &forbidden))
	s.userRepository.AssertExpectations(s.T())
}

func (s *UpgradeRequestServiceTestSuite) TestSubmitUpgradeRequest_PendingAlreadyExists() {
	user := &entity.User{Model: dbmodel.Model{ID: 1}, Type: enum.UserTypeCustomer}
	pending := &entity.UpgradeRequest{Model: dbmodel.Model{ID: 10}}
	req := upgraderequestdto.SubmitUpgradeRequestRequest{
		UserID:        1,
		RequestedType: enum.UserTypeFellow,
	}
	s.userRepository.On("FindUserByID", s.db, uint(1)).Return(user, nil).Once()
	s.upgradeRequestRepository.On("FindPendingByUserID", s.db, uint(1)).Return(pending, nil).Once()

	err := s.service.SubmitUpgradeRequest(req)

	s.Error(err)
	var forbidden exception.ForbiddenError
	s.True(errors.As(err, &forbidden))
	s.upgradeRequestRepository.AssertExpectations(s.T())
}

func (s *UpgradeRequestServiceTestSuite) TestSubmitUpgradeRequest_Success() {
	user := &entity.User{Model: dbmodel.Model{ID: 1}, Type: enum.UserTypeCustomer}
	var nilPending *entity.UpgradeRequest
	created := &entity.UpgradeRequest{Model: dbmodel.Model{ID: 11}}
	req := upgraderequestdto.SubmitUpgradeRequestRequest{
		UserID:        1,
		RequestedType: enum.UserTypeFellow,
		BusinessName:  "My Shop",
	}
	s.userRepository.On("FindUserByID", s.db, uint(1)).Return(user, nil).Once()
	s.upgradeRequestRepository.On("FindPendingByUserID", s.db, uint(1)).Return(nilPending, nil).Once()
	s.upgradeRequestRepository.On("CreateUpgradeRequest", s.db, entity.UpgradeRequest{
		UserID:        1,
		RequestedType: enum.UserTypeFellow,
		BusinessName:  "My Shop",
		Status:        enum.UpgradeRequestStatusPending,
	}).Return(created, nil).Once()

	err := s.service.SubmitUpgradeRequest(req)

	s.NoError(err)
	s.upgradeRequestRepository.AssertExpectations(s.T())
}

func (s *UpgradeRequestServiceTestSuite) TestGetAllUpgradeRequests_DefaultReturnsAll() {
	reqs := []*entity.UpgradeRequest{
		{Model: dbmodel.Model{ID: 1}, RequestedType: enum.UserTypeFellow, Status: enum.UpgradeRequestStatusPending},
	}
	s.upgradeRequestRepository.On("GetUpgradeRequests", s.db).Return(reqs, nil).Once()

	result, err := s.service.GetAllUpgradeRequests("")

	s.NoError(err)
	s.Len(result, 1)
	s.upgradeRequestRepository.AssertExpectations(s.T())
}

func (s *UpgradeRequestServiceTestSuite) TestGetAllUpgradeRequests_PendingFilter() {
	reqs := []*entity.UpgradeRequest{
		{Model: dbmodel.Model{ID: 2}, RequestedType: enum.UserTypeShopkeeperCash, Status: enum.UpgradeRequestStatusPending},
	}
	s.upgradeRequestRepository.On("GetUpgradeRequestsByStatus", s.db, enum.UpgradeRequestStatusPending).Return(reqs, nil).Once()

	result, err := s.service.GetAllUpgradeRequests("pending")

	s.NoError(err)
	s.Len(result, 1)
	s.upgradeRequestRepository.AssertExpectations(s.T())
}

func (s *UpgradeRequestServiceTestSuite) TestGetUpgradeRequest_NotFound() {
	var nilReq *entity.UpgradeRequest
	s.upgradeRequestRepository.On("FindUpgradeRequestByID", s.db, uint(99)).Return(nilReq, nil).Once()

	result, err := s.service.GetUpgradeRequest(99)

	s.Nil(result)
	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.upgradeRequestRepository.AssertExpectations(s.T())
}

func (s *UpgradeRequestServiceTestSuite) TestGetUpgradeRequest_RepoError() {
	repoErr := errors.New("db error")
	s.upgradeRequestRepository.On("FindUpgradeRequestByID", s.db, uint(1)).Return((*entity.UpgradeRequest)(nil), repoErr).Once()

	result, err := s.service.GetUpgradeRequest(1)

	s.Nil(result)
	s.ErrorIs(err, repoErr)
	s.upgradeRequestRepository.AssertExpectations(s.T())
}

func TestUpgradeRequestServiceSuite(t *testing.T) {
	suite.Run(t, new(UpgradeRequestServiceTestSuite))
}
