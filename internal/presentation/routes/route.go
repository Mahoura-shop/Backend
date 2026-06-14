package routes

import (
	"os"

	"github.com/Mahoura-shop/Backend/internal/presentation/middleware"
	httpv1 "github.com/Mahoura-shop/Backend/internal/presentation/routes/http/v1"
	"github.com/Mahoura-shop/Backend/wire"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(ginEngine *gin.Engine, app *wire.Application) {
	ginEngine.Use(app.Middlewares.CORS.CORS())
	ginEngine.Use(middleware.PrometheusMiddleware())
	ginEngine.Use(app.Middlewares.Logger.GinLoggerMiddleware)
	ginEngine.Use(app.Middlewares.Recovery.Recovery)

	ginEngine.GET("/metrics", gin.WrapH(promhttp.Handler()))
	ginEngine.GET("/health/ready", app.Controllers.General.HealthController.Ready)

	if os.Getenv("APP_MODE") == "test" {
		testGroup := ginEngine.Group("/v1/test")
		{
			testGroup.POST("", app.Controllers.General.TestController.Test)
			testGroup.POST("/reset", app.Controllers.General.TestController.Reset)
		}
	}
	ginEngine.Use(app.Middlewares.Localization.Localization)
	ginEngine.Use(app.Middlewares.RateLimit.RateLimit)

	ginEngine.OPTIONS("/*any", func(c *gin.Context) {})

	ginEngine.GET("/", func(c *gin.Context) {
		c.String(200, "backend is running")
	})

	v1 := ginEngine.Group("/v1")
	registerGeneralRoutes(v1, app)
	registerCustomerRoutes(v1, app)
	registerAdminRoutes(v1, app)
}

func registerGeneralRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	httpv1.SetupGeneralRoutes(v1, app)
}

func registerCustomerRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	customer := v1.Group("")
	customer.Use(app.Middlewares.Authentication.AuthRequired)
	httpv1.SetupCustomerRoutes(customer, app)
}

func registerAdminRoutes(v1 *gin.RouterGroup, app *wire.Application) {
	admin := v1.Group("")
	admin.Use(app.Middlewares.Authentication.AdminRequired)
	admin.Use(app.Middlewares.AdminActivityLog.LogActivity)
	httpv1.SetupAdminRoutes(admin, app)
}
