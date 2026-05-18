package user

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

const testUserID = uint(42)

type CustomerUserTestSuite struct {
	suite.Suite
	constants   *bootstrap.Constants
	userService *mocks.UserServiceMock
	controller  *CustomerUserController
	router      *gin.Engine
}

func (s *CustomerUserTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.userService = mocks.NewUserServiceMock()
	s.controller = NewCustomerUserController(s.constants, s.userService)

	s.router = customerUserRouter(s.constants, s.controller, false)
}

// customerUserRouter wires routes. withRecovery=true adds panic→HTTP conversion.
func customerUserRouter(c *bootstrap.Constants, ctrl *CustomerUserController, withRecovery bool) *gin.Engine {
	r := gin.New()
	r.Use(func(ctx *gin.Context) {
		ctx.Set(c.Context.Translator, mocks.NewTranslatorStub())
		ctx.Set(c.Context.ID, testUserID)
		ctx.Next()
	})
	if withRecovery {
		r.Use(panicToHTTP())
	}
	r.GET("/profile", ctrl.GetMyProfile)
	r.PATCH("/profile", ctrl.UpdateMyProfile)
	r.GET("/wallet", ctrl.GetUserWalletBalance)
	r.POST("/wallet/deposit", ctrl.DepositWallet)
	r.POST("/wallet/withdraw", ctrl.WithdrawWallet)
	return r
}

func panicToHTTP() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case *exception.AuthError:
					c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
				case exception.ValidationErrors:
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "validation"})
				case exception.BindingError:
					c.JSON(http.StatusBadRequest, gin.H{"error": "binding"})
				case exception.ConflictErrors:
					c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
				case exception.NotFoundError:
					c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
				}
				c.Abort()
			}
		}()
		c.Next()
	}
}

func jsonBuf(v any) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

// ─── Profile ──────────────────────────────────────────────────────────────────

// C-01: GET /profile → 200 with user data including type=regular
func (s *CustomerUserTestSuite) TestC01_GetProfile_Returns200() {
	s.userService.On("GetUserCredential", testUserID).Return(userdto.UserCredential{
		ID:    testUserID,
		Phone: "+989164911318",
		Type:  "regular",
	}, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/profile", nil)
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	s.Equal("regular", data["type"])
	s.userService.AssertExpectations(s.T())
}

// C-02: PATCH /profile → 200 with updated name
func (s *CustomerUserTestSuite) TestC02_UpdateProfile_Returns200() {
	req := userdto.UpdateProfileRequest{
		UserID:    testUserID,
		FirstName: "Ali",
		LastName:  "Rezaei",
		Email:     "",
	}
	s.userService.On("UpdateProfile", req).Return(userdto.UserCredential{
		ID:        testUserID,
		FirstName: "Ali",
		LastName:  "Rezaei",
		Type:      "regular",
	}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/profile",
		jsonBuf(map[string]string{"firstName": "Ali", "lastName": "Rezaei"}))
	r.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	data := body["data"].(map[string]any)
	s.Equal("Ali", data["firstName"])
	s.userService.AssertExpectations(s.T())
}

// C-03: PATCH /profile with invalid email → 422 (validation rejects it)
func (s *CustomerUserTestSuite) TestC03_UpdateProfile_InvalidEmail_Returns422() {
	router := customerUserRouter(s.constants, s.controller, true)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPatch, "/profile",
		jsonBuf(map[string]string{"email": "not-an-email"}))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
	s.userService.AssertNotCalled(s.T(), "UpdateProfile")
}

// C-04: Pricing check — GetProfile returns type=regular so UI shows Step4Price.
// Pricing resolution is tested in pricing_test.go; here we verify type field is correct.
func (s *CustomerUserTestSuite) TestC04_Profile_TypeIsRegular() {
	s.userService.On("GetUserCredential", testUserID).Return(userdto.UserCredential{
		ID:   testUserID,
		Type: "regular",
	}, nil)

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/profile", nil))

	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	s.Equal("regular", body["data"].(map[string]any)["type"])
}

// ─── Wallet ───────────────────────────────────────────────────────────────────

// C-22: GET /wallet → 200 with balance
func (s *CustomerUserTestSuite) TestC22_GetWallet_Returns200() {
	s.userService.On("GetUserWalletBalance", testUserID).Return(userdto.UserWalletBalance{Balance: 500_000}, nil)

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/wallet", nil))

	s.Equal(http.StatusOK, w.Code)
	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	s.Equal(float64(500_000), body["data"].(map[string]any)["balance"])
	s.userService.AssertExpectations(s.T())
}

// C-23: POST /wallet/deposit → 200, balance increases
func (s *CustomerUserTestSuite) TestC23_DepositWallet_Returns200() {
	s.userService.On("DepositWallet", userdto.UserBalanceUpdate{
		UserID: testUserID,
		Amount: 100_000,
	}).Return(userdto.UserWalletBalance{Balance: 600_000}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/wallet/deposit", jsonBuf(map[string]uint{"amount": 100_000}))
	r.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.userService.AssertExpectations(s.T())
}

// C-24: POST /wallet/withdraw sufficient balance → 200
func (s *CustomerUserTestSuite) TestC24_WithdrawWallet_SufficientBalance_Returns200() {
	s.userService.On("WithdrawWallet", userdto.UserBalanceUpdate{
		UserID: testUserID,
		Amount: 50_000,
	}).Return(userdto.UserWalletBalance{Balance: 450_000}, nil)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/wallet/withdraw", jsonBuf(map[string]uint{"amount": 50_000}))
	r.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, r)

	s.Equal(http.StatusOK, w.Code)
	s.userService.AssertExpectations(s.T())
}

// C-25: POST /wallet/withdraw insufficient balance → service returns error → 500
// (business rule enforced in service; controller just panics whatever service returns)
func (s *CustomerUserTestSuite) TestC25_WithdrawWallet_InsufficientBalance_ReturnsError() {
	s.userService.On("WithdrawWallet", userdto.UserBalanceUpdate{
		UserID: testUserID,
		Amount: 999_999_999,
	}).Return(userdto.UserWalletBalance{}, exception.ConflictErrors{
		Errors: []exception.FieldError{{Field: "wallet", Tag: "insufficientBalance"}},
	})

	router := customerUserRouter(s.constants, s.controller, true)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/wallet/withdraw", jsonBuf(map[string]uint{"amount": 999_999_999}))
	r.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, r)

	s.NotEqual(http.StatusOK, w.Code)
	s.userService.AssertExpectations(s.T())
}

func TestCustomerUserSuite(t *testing.T) {
	suite.Run(t, new(CustomerUserTestSuite))
}
