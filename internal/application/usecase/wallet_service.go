package usecase

import (
	walletdto "github.com/Mahoura-shop/Backend/internal/application/dto/wallet"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
)

type WalletService interface {
	ParseWallet(entity.Wallet) (walletdto.WalletCredential)
	GetWalletBalance(uint) (uint, error)
}
