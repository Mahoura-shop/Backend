package userdto

type AuthRequest struct {
	Phone string
}

type VerifyAuthRequest struct {
	Phone string
	OTP   string
}

type VerifyEmailRequest struct {
	UserID uint
	Email  string
	OTP    string
}

type AdminLoginRequest struct {
	Phone    string
	Password string
}

type UserBalanceUpdate struct {
	UserID uint
	Amount uint
}

type UpdateProductCountInCart struct {
	ProductID uint
	UserID    uint
}