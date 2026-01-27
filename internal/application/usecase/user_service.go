package usecase

import (
	userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
)

type UserService interface {
	IsUserActive(userID uint) error
	GetUserByID(userID uint) (*entity.User, error)
	GetUserCredential(userID uint) (userdto.CredentialResponse, error)
	BanUser(userID uint) error
	UnbanUser(userID uint) error
	Register(registerInfo userdto.BasicRegisterRequest) error
	VerifyPhone(verifyInfo userdto.VerifyPhoneRequest) error
	Login(loginInfo userdto.LoginRequest) (userdto.UserInfoResponse, error)
	ForgotPassword(forgotPasswordInfo userdto.ForgotPasswordRequest) error
	VerifyOTP(verifyInfo userdto.VerifyPhoneRequest) (userdto.UserInfoResponse, error)
	CompleteRegister(completeRegisterInfo userdto.CompleteRegisterRequest) error
	VerifyEmail(verifyOTPInfo userdto.VerifyEmailRequest) error
	ResetPassword(resetPassInfo userdto.ResetPasswordRequest) error
	FindActiveUserByPhone(phone string) (*entity.User, error)
	UpdateProfile(profileInfo userdto.UpdateProfileRequest) error
}
