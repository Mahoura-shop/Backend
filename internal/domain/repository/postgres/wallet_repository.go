package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type WalletRepository interface {
	FindWalletByID(database.Database, uint) (*entity.Wallet, error)
	FindWalletByPhone(database.Database, string) (*entity.Wallet, error)
	CreateWallet(database.Database, *entity.Wallet) (error)
	DeleteWalletByPhone(database.Database, string) (error)
	UpdateWallet(database.Database, *entity.Wallet) (error)
}
