package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type NotificationRepository struct{}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

func (r *NotificationRepository) Create(db database.Database, n entity.Notification) (*entity.Notification, error) {
	result := db.GetDB().Create(&n)
	return &n, result.Error
}

func (r *NotificationRepository) GetByUserID(db database.Database, userID uint) ([]*entity.Notification, error) {
	var notifications []*entity.Notification
	result := db.GetDB().Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications)
	return notifications, result.Error
}

func (r *NotificationRepository) MarkAsRead(db database.Database, id, userID uint) error {
	return db.GetDB().Model(&entity.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func (r *NotificationRepository) MarkAllAsRead(db database.Database, userID uint) error {
	return db.GetDB().Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Update("is_read", true).Error
}

func (r *NotificationRepository) CountUnread(db database.Database, userID uint) (int64, error) {
	var count int64
	result := db.GetDB().Model(&entity.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count)
	return count, result.Error
}
