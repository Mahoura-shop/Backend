package usecase

import notificationdto "github.com/Mahoura-shop/Backend/internal/application/dto/notification"

type NotificationService interface {
	CreateNotification(userID uint, notifType uint, title, body string, refID *uint) error
	GetNotifications(userID uint) ([]notificationdto.NotificationCredential, error)
	MarkAsRead(id, userID uint) error
	MarkAllAsRead(userID uint) error
	CountUnread(userID uint) (int64, error)
}
