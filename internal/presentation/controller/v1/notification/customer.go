package notification

import (
	"fmt"
	"time"

	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/sse"
	"github.com/Mahoura-shop/Backend/internal/presentation/controller"
	"github.com/gin-gonic/gin"
)

type CustomerNotificationController struct {
	constants           *bootstrap.Constants
	notificationService usecase.NotificationService
	jwtService          usecase.JWTService
}

func NewCustomerNotificationController(
	constants *bootstrap.Constants,
	notificationService usecase.NotificationService,
	jwtService usecase.JWTService,
) *CustomerNotificationController {
	return &CustomerNotificationController{
		constants:           constants,
		notificationService: notificationService,
		jwtService:          jwtService,
	}
}

func (c *CustomerNotificationController) GetNotifications(ctx *gin.Context) {
	userID, _ := ctx.Get(c.constants.Context.ID)
	notifications, err := c.notificationService.GetNotifications(userID.(uint))
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", notifications)
}

func (c *CustomerNotificationController) MarkAsRead(ctx *gin.Context) {
	type params struct {
		ID uint `uri:"id" validate:"required"`
	}
	p := controller.Validated[params](ctx)
	userID, _ := ctx.Get(c.constants.Context.ID)
	if err := c.notificationService.MarkAsRead(p.ID, userID.(uint)); err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", nil)
}

func (c *CustomerNotificationController) MarkAllAsRead(ctx *gin.Context) {
	userID, _ := ctx.Get(c.constants.Context.ID)
	if err := c.notificationService.MarkAllAsRead(userID.(uint)); err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", nil)
}

func (c *CustomerNotificationController) CountUnread(ctx *gin.Context) {
	userID, _ := ctx.Get(c.constants.Context.ID)
	count, err := c.notificationService.CountUnread(userID.(uint))
	if err != nil {
		panic(err)
	}
	controller.Response(ctx, 200, "", gin.H{"count": count})
}

// Stream handles SSE. Auth via ?token= query param since EventSource can't send headers.
func (c *CustomerNotificationController) Stream(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == "" {
		ctx.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	claims, err := c.jwtService.ValidateToken(token)
	if err != nil {
		ctx.JSON(401, gin.H{"message": "unauthorized"})
		return
	}

	userID := uint(claims["sub"].(float64))

	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("X-Accel-Buffering", "no")
	ctx.Writer.WriteHeader(200)
	ctx.Writer.Flush()

	ch := sse.Global.Subscribe(userID)
	defer sse.Global.Unsubscribe(userID, ch)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	clientGone := ctx.Request.Context().Done()

	for {
		select {
		case <-clientGone:
			return
		case event, ok := <-ch:
			if !ok {
				return
			}
			fmt.Fprintf(ctx.Writer, "data: %s\n\n", event.JSON())
			ctx.Writer.Flush()
		case <-ticker.C:
			fmt.Fprintf(ctx.Writer, ": ping\n\n")
			ctx.Writer.Flush()
		}
	}
}
