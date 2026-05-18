package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type TransactionRepository struct{}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (repo *TransactionRepository) FindTransactionByID(db database.Database, id uint) (entity.Transaction, error) {
	var transaction entity.Transaction
	result := db.GetDB().Preload("Wallet").First(&transaction, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return entity.Transaction{}, nil
		}
		return entity.Transaction{}, result.Error
	}
	return transaction, nil
}

func (repo *TransactionRepository) CreateTransaction(db database.Database, transaction entity.Transaction) error {
	return db.GetDB().Create(&transaction).Error
}

func (repo *TransactionRepository) UpdateTransaction(db database.Database, transaction entity.Transaction) error {
	return db.GetDB().Save(&transaction).Error
}

func (repo *TransactionRepository) FindTransactionsByWalletID(db database.Database, walletID uint) ([]entity.Transaction, error) {
	var transactions []entity.Transaction
	result := db.GetDB().Where("wallet_id = ?", walletID).Order("created_at DESC").Find(&transactions)
	return transactions, result.Error
}