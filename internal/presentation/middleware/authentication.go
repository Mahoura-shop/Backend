package middleware

import (
	"strings"

	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	repository "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	constants      *bootstrap.Constants
	jwtService     usecase.JWTService
	userRepository repository.UserRepository
	db             database.Database
	rbacConfig     *bootstrap.RBAC
}

func NewAuthMiddleware(
	constants *bootstrap.Constants,
	jwtService usecase.JWTService,
	userRepository repository.UserRepository,
	db database.Database,
	rbacConfig *bootstrap.RBAC,
) *AuthMiddleware {
	return &AuthMiddleware{
		constants:      constants,
		jwtService:     jwtService,
		userRepository: userRepository,
		db:             db,
		rbacConfig:     rbacConfig,
	}
}

func (am *AuthMiddleware) AuthRequired(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		unauthorizedError := exception.NewUnauthorizedError("empty auth header", nil)
		panic(unauthorizedError)
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		unauthorizedError := exception.NewUnauthorizedError("invalid token format", nil)
		panic(unauthorizedError)
	}

	tokenString := parts[1]
	if tokenString == "" {
		unauthorizedError := exception.NewUnauthorizedError("empty token", nil)
		panic(unauthorizedError)
	}

	claims, err := am.jwtService.ValidateToken(tokenString)
	if err != nil {
		panic(err)
	}

	ctx.Set(am.constants.Context.ID, uint(claims["sub"].(float64)))

	ctx.Next()
}

func (am *AuthMiddleware) AdminRequired(ctx *gin.Context) {
	am.AuthRequired(ctx)

	userID, _ := ctx.Get(am.constants.Context.ID)
	user, err := am.userRepository.FindUserByID(am.db, userID.(uint))
	if err != nil {
		panic(err)
	}
	if user == nil || !user.IsAdmin {
		panic(exception.NewAdminRequiredError())
	}

	ctx.Next()
}

func (am *AuthMiddleware) RequirePermission(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !am.rbacConfig.UseRBAC {
			ctx.Next()
			return
		}

		userID, _ := ctx.Get(am.constants.Context.ID)

		var user entity.User
		if err := am.db.GetDB().Preload("Role.Permissions").First(&user, userID.(uint)).Error; err != nil {
			panic(err)
		}

		if !user.IsAdmin {
			panic(exception.NewAdminRequiredError())
		}

		if user.RoleID == nil {
			ctx.Next()
			return
		}

		if user.Role != nil {
			for _, p := range user.Role.Permissions {
				if p.Name == permission {
					ctx.Next()
					return
				}
			}
		}

		panic(exception.NewInsufficientPermissionsError(permission))
	}
}