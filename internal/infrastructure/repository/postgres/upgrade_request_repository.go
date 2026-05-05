package postgres

import (
	"errors"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type UpgradeRequestRepository struct{}

func NewUpgradeRequestRepository() *UpgradeRequestRepository {
	return &UpgradeRequestRepository{}
}

func (r *UpgradeRequestRepository) CreateUpgradeRequest(db database.Database, req entity.UpgradeRequest) (*entity.UpgradeRequest, error) {
	result := db.GetDB().Create(&req)
	return &req, result.Error
}

func (r *UpgradeRequestRepository) FindUpgradeRequestByID(db database.Database, id uint) (*entity.UpgradeRequest, error) {
	var req entity.UpgradeRequest
	result := db.GetDB().Preload("User").Preload("ReviewedBy").Where("id = ?", id).First(&req)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &req, nil
}

func (r *UpgradeRequestRepository) GetUpgradeRequests(db database.Database) ([]*entity.UpgradeRequest, error) {
	var reqs []*entity.UpgradeRequest
	result := db.GetDB().Preload("User").Order("created_at DESC").Find(&reqs)
	return reqs, result.Error
}

func (r *UpgradeRequestRepository) GetUpgradeRequestsByUserID(db database.Database, userID uint) ([]*entity.UpgradeRequest, error) {
	var reqs []*entity.UpgradeRequest
	result := db.GetDB().Where("user_id = ?", userID).Order("created_at DESC").Find(&reqs)
	return reqs, result.Error
}

func (r *UpgradeRequestRepository) GetUpgradeRequestsByStatus(db database.Database, status enum.UpgradeRequestStatus) ([]*entity.UpgradeRequest, error) {
	var reqs []*entity.UpgradeRequest
	result := db.GetDB().Preload("User").Where("status = ?", status).Order("created_at DESC").Find(&reqs)
	return reqs, result.Error
}

func (r *UpgradeRequestRepository) FindPendingByUserID(db database.Database, userID uint) (*entity.UpgradeRequest, error) {
	var req entity.UpgradeRequest
	result := db.GetDB().Where("user_id = ? AND status = ?", userID, enum.UpgradeRequestStatusPending).First(&req)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &req, nil
}

func (r *UpgradeRequestRepository) UpdateUpgradeRequest(db database.Database, req entity.UpgradeRequest) error {
	return db.GetDB().Save(&req).Error
}
