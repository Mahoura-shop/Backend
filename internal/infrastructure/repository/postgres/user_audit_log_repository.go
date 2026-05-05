package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type UserAuditLogRepository struct{}

func NewUserAuditLogRepository() *UserAuditLogRepository {
	return &UserAuditLogRepository{}
}

func (r *UserAuditLogRepository) CreateUserAuditLog(db database.Database, log entity.UserAuditLog) error {
	return db.GetDB().Create(&log).Error
}

func (r *UserAuditLogRepository) GetAuditLogsByUserID(db database.Database, userID uint) ([]*entity.UserAuditLog, error) {
	var logs []*entity.UserAuditLog
	result := db.GetDB().Where("user_id = ?", userID).Order("created_at DESC").Find(&logs)
	return logs, result.Error
}
