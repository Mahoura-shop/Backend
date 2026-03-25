package usecase

import (
	userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
)

type UserService interface {
	ParseUser(entity.User) (userdto.UserCredential)
	IsUserActive(uint) error
	GetUserByID(uint) (*entity.User, error)
	FindActiveUserByPhone(string) (*entity.User, error)
	GetUserCredential(uint) (userdto.UserCredential, error)
	BanUser(uint) error
	UnbanUser(uint) error
	Auth(userdto.AuthRequest) error
	VerifyAuth(userdto.VerifyAuthRequest) (userdto.UserInfoResponse, error)
	FindUserByPhone(phone string) (*entity.User, error)
	VerifyEmail(userdto.VerifyEmailRequest) error
	AdminLogin(userdto.AdminLoginRequest) (userdto.AdminInfoResponse, error)
	GetDashboard() (userdto.DashboardResponse, error)
	GetUserWalletBalance(uint) (userdto.UserWalletBalance, error)
	DepositWallet(userdto.UserBalanceUpdate) (userdto.UserWalletBalance, error)
	WithdrawWallet(userdto.UserBalanceUpdate) (userdto.UserWalletBalance, error)
	AddProductToCart(userdto.AddProductToCartRequest) (error)
}
