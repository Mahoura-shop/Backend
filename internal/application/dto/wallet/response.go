package walletdto

import userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"

type WalletCredential struct {
	ID      uint
	Balance uint
	UserID  uint
	User    userdto.UserCredential
}

type WalletBalance struct {
	Balance uint
}