package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type NotificationRepository interface {
	Create(db database.Database, n entity.Notification) (*entity.Notification, error)
	GetByUserID(db database.Database, userID uint) ([]*entity.Notification, error)
	MarkAsRead(db database.Database, id, userID uint) error
	MarkAllAsRead(db database.Database, userID uint) error
	CountUnread(db database.Database, userID uint) (int64, error)
}
