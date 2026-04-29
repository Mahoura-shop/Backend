package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type InstalmentRepository struct{}

func NewInstalmentRepository() *InstalmentRepository {
	return &InstalmentRepository{}
}

func (repo *InstalmentRepository) CreateInstalment(db database.Database, instalment entity.Instalment) error {
	return db.GetDB().Create(&instalment).Error
}

func (repo *InstalmentRepository) GetInstalmentsByOrderID(db database.Database, orderID uint) ([]*entity.Instalment, error) {
	var instalments []*entity.Instalment
	result := db.GetDB().Where("order_id = ?", orderID).Order("number ASC").Find(&instalments)
	if result.Error != nil {
		return nil, result.Error
	}
	return instalments, nil
}

func (repo *InstalmentRepository) UpdateInstalment(db database.Database, instalment entity.Instalment) error {
	return db.GetDB().Save(&instalment).Error
}
