package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type AdminActivityLogRepository struct{}

func NewAdminActivityLogRepository() *AdminActivityLogRepository {
	return &AdminActivityLogRepository{}
}

func (r *AdminActivityLogRepository) Create(db database.Database, log entity.AdminActivityLog) error {
	return db.GetDB().Create(&log).Error
}

func (r *AdminActivityLogRepository) GetAll(db database.Database) ([]*entity.AdminActivityLog, error) {
	var logs []*entity.AdminActivityLog
	if err := db.GetDB().Preload("Admin").Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
