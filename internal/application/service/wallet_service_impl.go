package service

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	walletdto "github.com/Mahoura-shop/Backend/internal/application/dto/wallet"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type WalletService struct {
	constants        *bootstrap.Constants
	walletRepository postgres.WalletRepository
	userService      usecase.UserService
	db               database.Database
}

type WalletServiceDeps struct {
	Constants        *bootstrap.Constants
	WalletRepository postgres.WalletRepository
	UserService      usecase.UserService
	DB               database.Database
}

func NewWalletService(deps WalletServiceDeps) *WalletService {
	return &WalletService{
		constants:        deps.Constants,
		walletRepository: deps.WalletRepository,
		userService:      deps.UserService,
		db:               deps.DB,
	}
}

func (walletService *WalletService) ParseWallet(wallet entity.Wallet) (walletdto.WalletCredential) {
	response := walletdto.WalletCredential{
		ID:      wallet.ID,
		Balance: wallet.Balance,
		User:    walletService.userService.ParseUser(wallet.User),
	}

	return response
}