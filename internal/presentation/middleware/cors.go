package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type CORSMiddleware struct{}

func NewCorsMiddleware() *CORSMiddleware {
	return &CORSMiddleware{}
}

func (cm *CORSMiddleware) CORS() gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://127.0.0.1:3001", "http://127.0.0.1:3001"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "ngrok-skip-browser-warning"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	corsFn := cors.New(corsConfig)

	return func(c *gin.Context) {
		corsFn(c)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
		}
	}
}
