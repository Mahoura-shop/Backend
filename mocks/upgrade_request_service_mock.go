package mocks

import (
	upgraderequestdto "github.com/Mahoura-shop/Backend/internal/application/dto/upgrade_request"
	"github.com/stretchr/testify/mock"
)

type UpgradeRequestServiceMock struct {
	mock.Mock
}

func NewUpgradeRequestServiceMock() *UpgradeRequestServiceMock {
	return &UpgradeRequestServiceMock{}
}

func (m *UpgradeRequestServiceMock) SubmitUpgradeRequest(req upgraderequestdto.SubmitUpgradeRequestRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *UpgradeRequestServiceMock) GetMyUpgradeRequests(userID uint) ([]upgraderequestdto.UpgradeRequestCredential, error) {
	args := m.Called(userID)
	return args.Get(0).([]upgraderequestdto.UpgradeRequestCredential), args.Error(1)
}

func (m *UpgradeRequestServiceMock) GetAllUpgradeRequests(status string) ([]upgraderequestdto.UpgradeRequestCredential, error) {
	args := m.Called(status)
	return args.Get(0).([]upgraderequestdto.UpgradeRequestCredential), args.Error(1)
}

func (m *UpgradeRequestServiceMock) GetUpgradeRequest(requestID uint) (*upgraderequestdto.UpgradeRequestCredential, error) {
	args := m.Called(requestID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*upgraderequestdto.UpgradeRequestCredential), args.Error(1)
}

func (m *UpgradeRequestServiceMock) ReviewUpgradeRequest(requestID uint, req upgraderequestdto.ReviewUpgradeRequestRequest) error {
	args := m.Called(requestID, req)
	return args.Error(0)
}

func (m *UpgradeRequestServiceMock) ChangeUserType(userID uint, req upgraderequestdto.ChangeUserTypeRequest) error {
	args := m.Called(userID, req)
	return args.Error(0)
}

func (m *UpgradeRequestServiceMock) GetUserAuditLogs(userID uint) ([]upgraderequestdto.UserAuditLogCredential, error) {
	args := m.Called(userID)
	return args.Get(0).([]upgraderequestdto.UserAuditLogCredential), args.Error(1)
}
