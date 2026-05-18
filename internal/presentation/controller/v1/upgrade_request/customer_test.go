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

func init() { gin.SetMode(gin.TestMode) }

const upgradeUserID = uint(42)

type CustomerUpgradeRequestTestSuite struct {
	suite.Suite
	constants              *bootstrap.Constants
	upgradeRequestService  *mocks.UpgradeRequestServiceMock
	controller             *CustomerUpgradeRequestController
}

func (s *CustomerUpgradeRequestTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.upgradeRequestService = mocks.NewUpgradeRequestServiceMock()
	s.controller = NewCustomerUpgradeRequestController(s.constants, s.upgradeRequestService)
}

func (s *CustomerUpgradeRequestTestSuite) newRouter(withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set(s.constants.Context.Translator, mocks.NewTranslatorStub())
		c.Set(s.constants.Context.ID, upgradeUserID)
		c.Next()
	})
	if withRecovery {
		r.Use(upgradePanicToHTTP())
	}
	r.POST("/upgrade-requests", s.controller.SubmitUpgradeRequest)
	r.GET("/upgrade-requests", s.controller.GetMyUpgradeRequests)
	return r
}

func upgradePanicToHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case exception.ConflictErrors:
					c.JSON(http.StatusConflict, gin.H{"error": e.Error()})
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

// C-28: POST /upgrade-requests → 201, request created with status=pending
func (s *CustomerUpgradeRequestTestSuite) TestC28_SubmitUpgradeRequest_Returns201() {
	req := upgraderequestdto.SubmitUpgradeRequestRequest{
		UserID:        upgradeUserID,
		RequestedType: enum.UserTypeShopkeeperCash,
		BusinessName:  "My Shop",
		TaxID:         "1234567890",
	}
	s.upgradeRequestService.On("SubmitUpgradeRequest", req).Return(nil)

	body, _ := json.Marshal(map[string]any{
		"requestedType": enum.UserTypeShopkeeperCash,
		"businessName":  "My Shop",
		"taxID":         "1234567890",
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/upgrade-requests", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(false).ServeHTTP(w, r)

	s.Equal(http.StatusCreated, w.Code)
	s.upgradeRequestService.AssertExpectations(s.T())
}

// C-29: Duplicate request (pending exists) → service returns ConflictError → 409
func (s *CustomerUpgradeRequestTestSuite) TestC29_SubmitUpgradeRequest_Duplicate_Returns409() {
	req := upgraderequestdto.SubmitUpgradeRequestRequest{
		UserID:        upgradeUserID,
		RequestedType: enum.UserTypeShopkeeperCash,
		BusinessName:  "My Shop",
		TaxID:         "1234567890",
	}
	s.upgradeRequestService.On("SubmitUpgradeRequest", req).Return(
		exception.ConflictErrors{Errors: []exception.FieldError{{Field: "upgradeRequest", Tag: "alreadyPending"}}},
	)

	body, _ := json.Marshal(map[string]any{
		"requestedType": enum.UserTypeShopkeeperCash,
		"businessName":  "My Shop",
		"taxID":         "1234567890",
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/upgrade-requests", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(true).ServeHTTP(w, r)

	s.Equal(http.StatusConflict, w.Code)
	s.upgradeRequestService.AssertExpectations(s.T())
}

// C-29b: Missing required businessName → 422
func (s *CustomerUpgradeRequestTestSuite) TestC29b_SubmitUpgradeRequest_MissingBusinessName_Returns422() {
	body, _ := json.Marshal(map[string]any{
		"requestedType": enum.UserTypeShopkeeperCash,
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/upgrade-requests", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	s.newRouter(true).ServeHTTP(w, r)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
	s.upgradeRequestService.AssertNotCalled(s.T(), "SubmitUpgradeRequest")
}

// C-30: GET /upgrade-requests → 200 with own requests
func (s *CustomerUpgradeRequestTestSuite) TestC30_GetMyUpgradeRequests_Returns200() {
	s.upgradeRequestService.On("GetMyUpgradeRequests", upgradeUserID).Return(
		[]upgraderequestdto.UpgradeRequestCredential{
			{ID: 1, UserID: upgradeUserID, Status: enum.UpgradeRequestStatusPending},
		}, nil,
	)

	w := httptest.NewRecorder()
	s.newRouter(false).ServeHTTP(w,
		httptest.NewRequest(http.MethodGet, "/upgrade-requests", nil))

	s.Equal(http.StatusOK, w.Code)
	var resp map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &resp))
	requests := resp["data"].([]any)
	s.Len(requests, 1)
	s.upgradeRequestService.AssertExpectations(s.T())
}

func TestCustomerUpgradeRequestSuite(t *testing.T) {
	suite.Run(t, new(CustomerUpgradeRequestTestSuite))
}
