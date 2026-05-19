package notificationdto

import (
	"time"

	"github.com/Mahoura-shop/Backend/internal/domain/enum"
)

type NotificationCredential struct {
	ID        uint                 `json:"id"`
	Type      enum.NotificationType `json:"type"`
	Title     string               `json:"title"`
	Body      string               `json:"body"`
	IsRead    bool                 `json:"isRead"`
	RefID     *uint                `json:"refID,omitempty"`
	CreatedAt time.Time            `json:"createdAt"`
}
