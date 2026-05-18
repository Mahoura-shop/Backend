package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type AuthMiddlewareTestSuite struct {
	suite.Suite
	constants  *bootstrap.Constants
	jwtMock    *mocks.JwtServiceMock
	userRepo   *mocks.UserRepositoryMock
	dbMock     *mocks.DatabaseMock
	middleware *AuthMiddleware
}

func (s *AuthMiddlewareTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.jwtMock = mocks.NewJwtServiceMock()
	s.userRepo = mocks.NewUserRepositoryMock()
	s.dbMock = mocks.NewDatabaseMock()

	s.middleware = NewAuthMiddleware(
		s.constants,
		s.jwtMock,
		s.userRepo,
		s.dbMock,
	)
}

func newGinContext(method, path string) (*gin.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(method, path, nil)
	return c, w
}

// G-03: POST /cart/:id/add without token → AuthRequired panics with UnauthorizedError
// G-04: GET /cart without token → same middleware behaviour
// G-05: GET /profile without token → same middleware behaviour
// All three collapse to the same middleware contract: no Authorization header → panic(AuthError{Unauthorized})

func (s *AuthMiddlewareTestSuite) TestG03G04G05_NoAuthHeader_PanicsUnauthorized() {
	paths := []struct {
		id     string
		method string
		path   string
	}{
		{"G-03", "POST", "/cart/1/add"},
		{"G-04", "GET", "/cart"},
		{"G-05", "GET", "/profile"},
	}

	for _, tc := range paths {
		s.Run(tc.id+" "+tc.path, func() {
			c, _ := newGinContext(tc.method, tc.path)
			// No Authorization header set

			s.Panics(func() {
				s.middleware.AuthRequired(c)
			}, "AuthRequired must panic when Authorization header is missing")

			// Verify the panic is specifically an AuthError (not a generic panic)
			var caughtErr any
			func() {
				defer func() { caughtErr = recover() }()
				s.middleware.AuthRequired(c)
			}()

			authErr, ok := caughtErr.(*exception.AuthError)
			assert.True(s.T(), ok, "%s: panic value must be *exception.AuthError, got %T", tc.id, caughtErr)
			if ok {
				assert.Equal(s.T(), exception.ErrorTypeUnauthorized, authErr.Type,
					"%s: AuthError.Type must be ErrorTypeUnauthorized", tc.id)
			}
		})
	}
}

// Extra: malformed Authorization header (no "Bearer" prefix) → UnauthorizedError
func (s *AuthMiddlewareTestSuite) TestAuthRequired_MalformedHeader_PanicsUnauthorized() {
	c, _ := newGinContext("GET", "/profile")
	c.Request.Header.Set("Authorization", "InvalidToken abc123")

	var caughtErr any
	func() {
		defer func() { caughtErr = recover() }()
		s.middleware.AuthRequired(c)
	}()

	authErr, ok := caughtErr.(*exception.AuthError)
	s.True(ok, "panic must be *exception.AuthError, got %T", caughtErr)
	if ok {
		s.Equal(exception.ErrorTypeUnauthorized, authErr.Type)
	}
}

// Extra: empty token string after "Bearer " → UnauthorizedError
func (s *AuthMiddlewareTestSuite) TestAuthRequired_EmptyBearerToken_PanicsUnauthorized() {
	c, _ := newGinContext("GET", "/profile")
	c.Request.Header.Set("Authorization", "Bearer ")

	var caughtErr any
	func() {
		defer func() { caughtErr = recover() }()
		s.middleware.AuthRequired(c)
	}()

	authErr, ok := caughtErr.(*exception.AuthError)
	s.True(ok, "panic must be *exception.AuthError, got %T", caughtErr)
	if ok {
		s.Equal(exception.ErrorTypeUnauthorized, authErr.Type)
	}
}

// ─── AdminRequired ────────────────────────────────────────────────────────────

// AA-04a: no token on admin route → panic(AuthError{Unauthorized})
func (s *AuthMiddlewareTestSuite) TestAA04a_AdminRequired_NoToken_PanicsUnauthorized() {
	c, _ := newGinContext("GET", "/admin/dashboard")

	var caughtErr any
	func() {
		defer func() { caughtErr = recover() }()
		s.middleware.AdminRequired(c)
	}()

	authErr, ok := caughtErr.(*exception.AuthError)
	s.True(ok, "panic must be *exception.AuthError, got %T", caughtErr)
	if ok {
		s.Equal(exception.ErrorTypeUnauthorized, authErr.Type)
	}
}

// AA-04b: valid JWT but non-admin user → panic(ForbiddenError)
func (s *AuthMiddlewareTestSuite) TestAA04b_AdminRequired_NonAdminUser_PanicsForbidden() {
	// jwtService returns valid claims with sub=1
	s.jwtMock.On("ValidateToken", "validtoken").
		Return(map[string]any{"sub": float64(1)}, nil)

	// userRepository returns a user that is NOT admin
	nonAdmin := &entity.User{IsAdmin: false}
	s.userRepo.On("FindUserByID", s.dbMock, uint(1)).Return(nonAdmin, nil)

	c, _ := newGinContext("GET", "/admin/dashboard")
	c.Request.Header.Set("Authorization", "Bearer validtoken")

	var caughtErr any
	func() {
		defer func() { caughtErr = recover() }()
		s.middleware.AdminRequired(c)
	}()

	_, isForbidden := caughtErr.(exception.ForbiddenError)
	_, isAuthError := caughtErr.(*exception.AuthError)
	s.True(isForbidden || isAuthError,
		"non-admin user must trigger ForbiddenError or AuthError, got %T", caughtErr)

	s.jwtMock.AssertExpectations(s.T())
	s.userRepo.AssertExpectations(s.T())
}

// AA-04c: valid admin JWT → AdminRequired does NOT panic
func (s *AuthMiddlewareTestSuite) TestAA04c_AdminRequired_AdminUser_PassesThrough() {
	s.jwtMock.On("ValidateToken", "admintoken").
		Return(map[string]any{"sub": float64(99)}, nil)

	adminUser := &entity.User{IsAdmin: true}
	s.userRepo.On("FindUserByID", s.dbMock, uint(99)).Return(adminUser, nil)

	c, _ := newGinContext("GET", "/admin/dashboard")
	c.Request.Header.Set("Authorization", "Bearer admintoken")

	// gin.Context.Next() panics when there are no handlers in test context;
	// use a real router to avoid that.
	router := gin.New()
	router.GET("/admin/dashboard", func(ctx *gin.Context) {
		s.middleware.AdminRequired(ctx)
		ctx.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer admintoken")
	router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)
	s.jwtMock.AssertExpectations(s.T())
	s.userRepo.AssertExpectations(s.T())
}

func TestAuthMiddlewareSuite(t *testing.T) {
	suite.Run(t, new(AuthMiddlewareTestSuite))
}
