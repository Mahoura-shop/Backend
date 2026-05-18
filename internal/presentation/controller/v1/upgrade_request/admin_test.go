package upgraderequest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	upgraderequestdto "github.com/Mahoura-shop/Backend/internal/application/dto/upgrade_request"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

const adminReqID = uint(20)
const adminID = uint(1)
const upgradeTargetUserID = uint(42)

type AdminUpgradeRequestTestSuite struct {
	suite.Suite
	constants             *bootstrap.Constants
	upgradeRequestService *mocks.UpgradeRequestServiceMock
	controller            *AdminUpgradeRequestController
}

func (s *AdminUpgradeRequestTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.upgradeRequestService = mocks.NewUpgradeRequestServiceMock()
	s.controller = NewAdminUpgradeRequestController(s.constants, s.upgradeRequestService)
}

func (s *AdminUpgradeRequestTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Set(s.constants.Context.ID, adminID)
		c.Next()
	})
	if withRecovery {
		r.Use(adminUpgradePanicToHTTP())
	}
	r.GET("/admin/upgrade-requests", s.controller.GetUpgradeRequests)
	r.GET("/admin/upgrade-requests/:requestID", s.controller.GetUpgradeRequest)
	r.POST("/admin/upgrade-requests/:requestID/review", s.controller.ReviewUpgradeRequest)
	r.POST("/admin/users/:userID/change-type", s.controller.ChangeUserType)
	r.GET("/admin/users/:userID/audit-logs", s.controller.GetUserAuditLogs)
	return r
}

func adminUpgradePanicToHTTP() gin.HandlerFunc {
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

// ADM-40: GET /admin/upgrade-requests → 200 with all requests
func (s *AdminUpgradeRequestTestSuite) TestADM40_GetAllRequests_Returns200() {
	s.upgradeRequestService.On("GetAllUpgradeRequests", "").Return(
		[]upgraderequestdto.UpgradeRequestCredential{
			{ID: adminReqID, UserID: upgradeTargetUserID, Status: enum.UpgradeRequestStatusPending},
			{ID: 21, UserID: 43, Status: enum.UpgradeRequestStatusApproved},
		}, nil,
	)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/admin/upgrade-requests", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	requests := body["data"].([]any)
	s.Len(requests, 2)
	s.upgradeRequestService.AssertExpectations(s.T())
}

// ADM-41: GET /admin/upgrade-requests?status=pending → 200, filtered list
func (s *AdminUpgradeRequestTestSuite) TestADM41_GetRequests_FilterByStatus_Returns200() {
	s.upgradeRequestService.On("GetAllUpgradeRequests", "pending").Return(
		[]upgraderequestdto.UpgradeRequestCredential{
			{ID: adminReqID, Status: enum.UpgradeRequestStatusPending},
		}, nil,
	)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/admin/upgrade-requests?status=pending", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	requests := body["data"].([]any)
	s.Len(requests, 1)
	s.upgradeRequestService.AssertExpectations(s.T())
}

// ADM-42: GET /admin/upgrade-requests/:requestID → 200 with full request detail
func (s *AdminUpgradeRequestTestSuite) TestADM42_GetRequestDetail_Returns200() {
	s.upgradeRequestService.On("GetUpgradeRequest", adminReqID).Return(
		&upgraderequestdto.UpgradeRequestCredential{
			ID:           adminReqID,
			UserID:       upgradeTargetUserID,
			BusinessName: "My Shop",
			Status:       enum.UpgradeRequestStatusPending,
		}, nil,
	)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/admin/upgrade-requests/20", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	s.Equal(float64(adminReqID), data["id"])
	s.upgradeRequestService.AssertExpectations(s.T())
}

// ADM-43: POST /admin/upgrade-requests/:requestID/review action=approve → 200
func (s *AdminUpgradeRequestTestSuite) TestADM43_ApproveRequest_Returns200() {
	req := upgraderequestdto.ReviewUpgradeRequestRequest{
		Action:    "approve",
		AdminNote: "looks good",
		AdminID:   adminID,
	}
	s.upgradeRequestService.On("ReviewUpgradeRequest", adminReqID, req).Return(nil)

	payload, _ := json.Marshal(map[string]any{"action": "approve", "adminNote": "looks good"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/upgrade-requests/20/review", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.upgradeRequestService.AssertExpectations(s.T())
}

// ADM-44: POST /admin/upgrade-requests/:requestID/review action=reject → 200
func (s *AdminUpgradeRequestTestSuite) TestADM44_RejectRequest_Returns200() {
	req := upgraderequestdto.ReviewUpgradeRequestRequest{
		Action:  "reject",
		AdminID: adminID,
	}
	s.upgradeRequestService.On("ReviewUpgradeRequest", adminReqID, req).Return(nil)

	payload, _ := json.Marshal(map[string]any{"action": "reject"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/upgrade-requests/20/review", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.upgradeRequestService.AssertExpectations(s.T())
}

// ADM-45: POST /admin/upgrade-requests/:requestID/review action=info → 200
func (s *AdminUpgradeRequestTestSuite) TestADM45_RequestMoreInfo_Returns200() {
	req := upgraderequestdto.ReviewUpgradeRequestRequest{
		Action:    "info",
		AdminNote: "need tax docs",
		AdminID:   adminID,
	}
	s.upgradeRequestService.On("ReviewUpgradeRequest", adminReqID, req).Return(nil)

	payload, _ := json.Marshal(map[string]any{"action": "info", "adminNote": "need tax docs"})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/upgrade-requests/20/review", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.upgradeRequestService.AssertExpectations(s.T())
}

// ADM-21/22/23/24 — ChangeUserType: promote/demote via POST /admin/users/:userID/change-type
func (s *AdminUpgradeRequestTestSuite) TestADM21_PromoteToShopkeeperCash_Returns200() {
	req := upgraderequestdto.ChangeUserTypeRequest{
		NewType: uint(enum.UserTypeShopkeeperCash),
		AdminID: adminID,
	}
	s.upgradeRequestService.On("ChangeUserType", upgradeTargetUserID, req).Return(nil)

	payload, _ := json.Marshal(map[string]any{
		"newType": uint(enum.UserTypeShopkeeperCash),
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/users/42/change-type", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.upgradeRequestService.AssertExpectations(s.T())
}

func (s *AdminUpgradeRequestTestSuite) TestADM23_PromoteToFellow_Returns200() {
	req := upgraderequestdto.ChangeUserTypeRequest{
		NewType: uint(enum.UserTypeFellow),
		AdminID: adminID,
	}
	s.upgradeRequestService.On("ChangeUserType", upgradeTargetUserID, req).Return(nil)

	payload, _ := json.Marshal(map[string]any{
		"newType": uint(enum.UserTypeFellow),
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/users/42/change-type", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.upgradeRequestService.AssertExpectations(s.T())
}

func (s *AdminUpgradeRequestTestSuite) TestADM24_DemoteToRegular_Returns200() {
	req := upgraderequestdto.ChangeUserTypeRequest{
		NewType: uint(enum.UserTypeCustomer),
		AdminID: adminID,
	}
	s.upgradeRequestService.On("ChangeUserType", upgradeTargetUserID, req).Return(nil)

	payload, _ := json.Marshal(map[string]any{
		"newType": uint(enum.UserTypeCustomer),
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/admin/users/42/change-type", bytes.NewBuffer(payload))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.upgradeRequestService.AssertExpectations(s.T())
}

// ADM-25: GET /admin/users/:userID/audit-logs → 200 with type change history
func (s *AdminUpgradeRequestTestSuite) TestADM25_GetUserAuditLogs_Returns200() {
	s.upgradeRequestService.On("GetUserAuditLogs", upgradeTargetUserID).Return(
		[]upgraderequestdto.UserAuditLogCredential{
			{ID: 1, UserID: upgradeTargetUserID, OldType: "regular", NewType: "shopkeeperCash"},
		}, nil,
	)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/admin/users/42/audit-logs", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	logs := body["data"].([]any)
	s.Len(logs, 1)
	s.upgradeRequestService.AssertExpectations(s.T())
}

func TestAdminUpgradeRequestSuite(t *testing.T) {
	suite.Run(t, new(AdminUpgradeRequestTestSuite))
}
