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

func init() {
	gin.SetMode(gin.TestMode)
}

type AuthTestSuite struct {
	suite.Suite
	constants   *bootstrap.Constants
	userService *mocks.UserServiceMock
	jwtMock     *mocks.JwtServiceMock
	controller  *GeneralUserController
	router      *gin.Engine
}

func (s *AuthTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.userService = mocks.NewUserServiceMock()
	s.jwtMock = mocks.NewJwtServiceMock()

	s.controller = NewGeneralUserController(s.constants, s.userService, s.jwtMock)

	s.router = gin.New()

	// Inject translator stub so GetTranslator(ctx, ...) doesn't panic
	translatorKey := s.constants.Context.Translator
	s.router.Use(func(c *gin.Context) {
		c.Set(translatorKey, mocks.NewTranslatorStub())
		c.Next()
	})

	s.router.POST("/auth", s.controller.Auth)
	s.router.POST("/auth/verify", s.controller.VerifyAuth)
}

func jsonBody(v any) *bytes.Buffer {
	b, _ := json.Marshal(v)
	return bytes.NewBuffer(b)
}

// ─── OTP Flow ────────────────────────────────────────────────────────────────

// A-01: valid phone → OTP sent → 200
func (s *AuthTestSuite) TestA01_Auth_ValidPhone_Returns200() {
	s.userService.On("Auth", userdto.AuthRequest{Phone: "+989164911318"}).Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth",
		jsonBody(map[string]string{"phone": "+989164911318"}))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.userService.AssertExpectations(s.T())
}

// A-02: correct OTP → JWT returned → 200 with accessToken + refreshToken
func (s *AuthTestSuite) TestA02_VerifyAuth_CorrectOTP_ReturnsTokens() {
	resp := userdto.UserInfoResponse{
		AccessToken:  "access.token.here",
		RefreshToken: "refresh.token.here",
		Type:         "regular",
	}
	s.userService.On("VerifyAuth", userdto.VerifyAuthRequest{
		Phone: "+989164911318",
		OTP:   "123456",
	}).Return(resp, nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/verify",
		jsonBody(map[string]string{"phone": "+989164911318", "otp": "123456"}))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var body map[string]any
	s.Require().NoError(json.Unmarshal(w.Body.Bytes(), &body))
	data, ok := body["data"].(map[string]any)
	s.Require().True(ok, "data must be an object")
	s.NotEmpty(data["accessToken"], "accessToken must be present")
	s.NotEmpty(data["refreshToken"], "refreshToken must be present")
	s.userService.AssertExpectations(s.T())
}

// A-03: wrong OTP → service returns AuthError → 401
func (s *AuthTestSuite) TestA03_VerifyAuth_WrongOTP_Returns401() {
	s.userService.On("VerifyAuth", userdto.VerifyAuthRequest{
		Phone: "+989164911318",
		OTP:   "000000",
	}).Return(userdto.UserInfoResponse{}, exception.NewUnauthorizedError("invalid otp", nil))

	// Use a router with recovery so the panic is converted to an HTTP response
	router := s.routerWithRecovery()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/verify",
		jsonBody(map[string]string{"phone": "+989164911318", "otp": "000000"}))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
	s.userService.AssertExpectations(s.T())
}

// A-04: expired OTP → service returns ExpiredTokenError → 401
func (s *AuthTestSuite) TestA04_VerifyAuth_ExpiredOTP_Returns401() {
	s.userService.On("VerifyAuth", userdto.VerifyAuthRequest{
		Phone: "+989164911318",
		OTP:   "123456",
	}).Return(userdto.UserInfoResponse{}, exception.NewExpiredTokenError(nil))

	router := s.routerWithRecovery()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/verify",
		jsonBody(map[string]string{"phone": "+989164911318", "otp": "123456"}))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)
	s.userService.AssertExpectations(s.T())
}

// A-05: existing phone → Auth called → 200 (service handles existing user internally)
func (s *AuthTestSuite) TestA05_Auth_ExistingPhone_Returns200() {
	s.userService.On("Auth", userdto.AuthRequest{Phone: "+989164911318"}).Return(nil)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth",
		jsonBody(map[string]string{"phone": "+989164911318"}))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.userService.AssertExpectations(s.T())
}

// A-06: verify OTP with missing phone field → validation error → 422
// (no prior OTP request needed — controller validates input before calling service)
func (s *AuthTestSuite) TestA06_VerifyAuth_MissingPhone_Returns422() {
	router := s.routerWithRecovery()
	w := httptest.NewRecorder()
	// send OTP but no phone
	req := httptest.NewRequest(http.MethodPost, "/auth/verify",
		jsonBody(map[string]string{"otp": "123456"}))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	s.Equal(http.StatusUnprocessableEntity, w.Code)
	// service must NOT be called
	s.userService.AssertNotCalled(s.T(), "VerifyAuth")
}


// ─── helpers ─────────────────────────────────────────────────────────────────

// routerWithRecovery returns a router that converts panics to HTTP responses.
// Uses a minimal inline recovery so tests don't need the full middleware stack.
func (s *AuthTestSuite) routerWithRecovery() *gin.Engine {
	r := gin.New()

	translatorKey := s.constants.Context.Translator
	r.Use(func(c *gin.Context) {
		c.Set(translatorKey, mocks.NewTranslatorStub())
		c.Next()
	})

	r.Use(func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case *exception.AuthError:
					c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
				case exception.ValidationErrors:
					c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "validation"})
				case exception.BindingError:
					c.JSON(http.StatusBadRequest, gin.H{"error": "binding"})
				default:
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
				}
				c.Abort()
			}
		}()
		c.Next()
	})

	r.POST("/auth", s.controller.Auth)
	r.POST("/auth/verify", s.controller.VerifyAuth)
	return r
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}
