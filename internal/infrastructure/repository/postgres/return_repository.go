package postgres

import (
	"errors"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type ReturnRepository struct{}

func NewReturnRepository() *ReturnRepository {
	return &ReturnRepository{}
}

func (r *ReturnRepository) CreateReturn(db database.Database, ret entity.Return) (*entity.Return, error) {
	result := db.GetDB().Create(&ret)
	return &ret, result.Error
}

func (r *ReturnRepository) FindReturnByID(db database.Database, id uint) (*entity.Return, error) {
	var ret entity.Return
	result := db.GetDB().
		Preload("OrderItem.Product").
		Preload("User").
		Preload("ApprovedBy").
		Where("id = ?", id).First(&ret)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &ret, result.Error
}

func (r *ReturnRepository) GetReturnsByUserID(db database.Database, userID uint) ([]*entity.Return, error) {
	var rets []*entity.Return
	result := db.GetDB().
		Preload("OrderItem.Product").
		Where("user_id = ?", userID).
		Order("requested_at DESC").Find(&rets)
	return rets, result.Error
}

func (r *ReturnRepository) GetAllReturns(db database.Database) ([]*entity.Return, error) {
	var rets []*entity.Return
	result := db.GetDB().
		Preload("OrderItem.Product").
		Preload("User").
		Order("requested_at DESC").Find(&rets)
	return rets, result.Error
}

func (r *ReturnRepository) GetReturnsByStatus(db database.Database, status enum.ReturnStatus) ([]*entity.Return, error) {
	var rets []*entity.Return
	result := db.GetDB().
		Preload("OrderItem.Product").
		Preload("User").
		Where("status = ?", status).
		Order("requested_at DESC").Find(&rets)
	return rets, result.Error
}

func (r *ReturnRepository) UpdateReturn(db database.Database, ret entity.Return) error {
	return db.GetDB().Save(&ret).Error
}
