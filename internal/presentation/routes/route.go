package routes

import (
	httpv1 "github.com/Mahoura-shop/Backend/internal/presentation/routes/http/v1"
	"github.com/Mahoura-shop/Backend/wire"
	"github.com/gin-gonic/gin"
)

func Run(ginEngine *gin.Engine, app *wire.Application) {
	ginEngine.Use(app.Middlewares.CORS.CORS())
	ginEngine.Use(app.Middlewares.Logger.GinLoggerMiddleware)
	ginEngine.Use(app.Middlewares.Localization.Localization)
	ginEngine.Use(app.Middlewares.Recovery.Recovery)
	ginEngine.Use(app.Middlewares.RateLimit.RateLimit)

	v1 := ginEngine.Group("/v1")
	registerGeneralRoutes(v1, app)
	registerAdminRoutes(v1, app)
}

func registerGeneralRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	httpv1.SetupGeneralRoutes(v1, app)
}

func registerAdminRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	admin := v1.Group("")
	admin.Use(app.Middlewares.Authentication.AuthRequired)
	httpv1.SetupAdminRoutes(admin, app)
}
