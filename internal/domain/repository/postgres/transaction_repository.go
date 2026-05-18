package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type TransactionRepository interface {
	FindTransactionByID(database.Database, uint) (entity.Transaction, error)
	FindTransactionsByWalletID(database.Database, uint) ([]entity.Transaction, error)
	CreateTransaction(database.Database, entity.Transaction) (error)
	UpdateTransaction(database.Database, entity.Transaction) (error)
}
