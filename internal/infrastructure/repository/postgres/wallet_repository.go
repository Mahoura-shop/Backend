package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type WalletRepository struct{}

func NewWalletRepository() *WalletRepository {
	return &WalletRepository{}
}

func (repo *WalletRepository) FindWalletByID(db database.Database, id uint) (*entity.Wallet, error) {
	var wallet entity.Wallet
	result := db.GetDB().First(&wallet, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &wallet, nil
}

func (repo *WalletRepository) FindWalletByPhone(db database.Database, phone string) (*entity.Wallet, error) {
	var wallet entity.Wallet
	result := db.GetDB().Where("phone = ?", phone).First(&wallet)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &wallet, nil
}

func (repo *WalletRepository) CreateWallet(db database.Database, wallet *entity.Wallet) error {
	return db.GetDB().Create(&wallet).Error
}

func (repo *WalletRepository) DeleteWalletByPhone(db database.Database, phone string) error {
	return db.GetDB().Where("phone = ?", phone).Unscoped().Delete(&entity.Wallet{}).Error
}

func (repo *WalletRepository) UpdateWallet(db database.Database, wallet *entity.Wallet) error {
	return db.GetDB().Save(&wallet).Error
}