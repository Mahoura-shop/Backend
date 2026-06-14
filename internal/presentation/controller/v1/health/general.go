package health

import (
	"context"
	"net/http"

	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/gin-gonic/gin"
)

type HealthController struct {
	db  database.Database
	rdb database.Cache
}

func NewHealthController(db database.Database, rdb database.Cache) *HealthController {
	return &HealthController{db: db, rdb: rdb}
}

func (h *HealthController) Ready(c *gin.Context) {
	deps := gin.H{}
	allOK := true

	sqlDB, err := h.db.GetDB().DB()
	if err != nil || sqlDB.PingContext(context.Background()) != nil {
		deps["postgres"] = "down"
		allOK = false
	} else {
		deps["postgres"] = "up"
	}

	if err := h.rdb.GetRDB().Ping(context.Background()).Err(); err != nil {
		deps["redis"] = "down"
		allOK = false
	} else {
		deps["redis"] = "up"
	}

	if allOK {
		c.JSON(http.StatusOK, gin.H{"status": "ready", "dependencies": deps})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready", "dependencies": deps})
	}
}
