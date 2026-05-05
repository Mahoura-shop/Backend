package service

import (
	"context"
	"time"

	"github.com/Mahoura-shop/Backend/bootstrap"
	userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"
	"github.com/Mahoura-shop/Backend/internal/application/usecase"
	"github.com/Mahoura-shop/Backend/internal/domain/communication"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/redis"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type UserService struct {
	constants             *bootstrap.Constants
	otpService            usecase.OTPService
	jwtService            usecase.JWTService
	smsService            communication.SMSService
	emailService          communication.EmailService
	userRepository        postgres.UserRepository
	categoryRepository    postgres.CategoryRepository
	brandRepository       postgres.BrandRepository
	productRepository     postgres.ProductRepository
	walletRepository      postgres.WalletRepository
	cartRepository        postgres.CartRepository
	transactionRepository postgres.TransactionRepository
	userCacheRepository   redis.UserCacheRepository
	db                    database.Database
}

type UserServiceDeps struct {
	Constants             *bootstrap.Constants
	OTPService            usecase.OTPService
	JWTService            usecase.JWTService
	SMSService            communication.SMSService
	EmailService          communication.EmailService
	UserRepository        postgres.UserRepository
	CategoryRepository    postgres.CategoryRepository
	BrandRepository       postgres.BrandRepository
	ProductRepository     postgres.ProductRepository
	WalletRepository      postgres.WalletRepository
	CartRepository        postgres.CartRepository
	TransactionRepository postgres.TransactionRepository
	UserCacheRepository   redis.UserCacheRepository
	DB                    database.Database
}

func NewUserService(deps UserServiceDeps) *UserService {
	return &UserService{
		constants:             deps.Constants,
		otpService:            deps.OTPService,
		jwtService:            deps.JWTService,
		smsService:            deps.SMSService,
		emailService:          deps.EmailService,
		userRepository:        deps.UserRepository,
		categoryRepository:    deps.CategoryRepository,
		brandRepository:       deps.BrandRepository,
		productRepository:     deps.ProductRepository,
		walletRepository:      deps.WalletRepository,
		cartRepository:        deps.CartRepository,
		transactionRepository: deps.TransactionRepository,
		userCacheRepository:   deps.UserCacheRepository,
		db:                    deps.DB,
	}
}


func (userService *UserService) ParseUser(user entity.User) (userdto.UserCredential) {
	response := userdto.UserCredential{
		ID:            user.ID,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Phone:         user.Phone,
		Email:         user.Email,
		ProfilePic:    user.ProfilePicPath,
		Status:        user.Status.String(),
		Type:          user.Type.String(),
	}

	return response
}

func (userService *UserService) IsUserActive(userID uint) error {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user.Status == enum.UserStatusBlock {
		return exception.NewBannedUserForbiddenError()
	}
	return nil
}

func (userService *UserService) GetUserByID(userID uint) (*entity.User, error) {
	user, err := userService.userRepository.FindUserByID(userService.db, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return nil, notFoundError
	}
	return user, nil
}

func (userService *UserService) FindActiveUserByPhone(phone string) (*entity.User, error) {
	user, err := userService.userRepository.FindUserByPhone(userService.db, phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return nil, notFoundError
	}

	return user, nil
}

func (userService *UserService) GetUserCredential(userID uint) (userdto.UserCredential, error) {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return userdto.UserCredential{}, err
	}

	// return userdto.UserCredential{
	// 	ID:         user.ID,
	// 	FirstName:  user.FirstName,
	// 	LastName:   user.LastName,
	// 	Phone:      user.Phone,
	// 	Email:      user.Email,
	// 	Status:     user.Status.String(),
	// 	Type:       user.Type.String(),
	// },
	return userService.ParseUser(*user), nil
}

func (userService *UserService) GetUsers() ([]userdto.UserCredential, error) {
	users, err := userService.userRepository.FindUsers(userService.db)
	if err != nil {
		return nil, err
	}
	var result []userdto.UserCredential
	for _, u := range users {
		result = append(result, userService.ParseUser(*u))
	}
	return result, nil
}

func (userService *UserService) BanUser(userID uint) error {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.Status == enum.UserStatusBlock {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.User, userService.constants.Tag.AlreadyBlocked)
		return conflictErrors
	}
	user.Status = enum.UserStatusBlock
	err = userService.userRepository.UpdateUser(userService.db, *user)
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) UnbanUser(userID uint) error {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.Status == enum.UserStatusActive {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.User, userService.constants.Tag.AlreadyActive)
		return conflictErrors
	}
	user.Status = enum.UserStatusActive
	err = userService.userRepository.UpdateUser(userService.db, *user)
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) Auth(authInfo userdto.AuthRequest) error {
	otp, expiryMinute, err := userService.otpService.GenerateOTP()
	if err != nil {
		return err
	}
	redisKey := userService.constants.RedisKey.GenerateOTPKey(authInfo.Phone)
	err = userService.userCacheRepository.Set(context.Background(), redisKey, otp, time.Duration(expiryMinute)*time.Minute)
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) VerifyAuth(verifyAuthInfo userdto.VerifyAuthRequest) (userdto.UserInfoResponse, error) {
	user, err := userService.userRepository.FindUserByPhone(userService.db, verifyAuthInfo.Phone)
	
	redisKey := userService.constants.RedisKey.GenerateOTPKey(verifyAuthInfo.Phone)
	err = userService.otpService.VerifyOTP(redisKey, verifyAuthInfo.OTP)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}
	
	if user == nil {
		err = userService.db.WithTransaction(func(tx database.Database) error {
			user := &entity.User{
				Phone:         verifyAuthInfo.Phone,
				EmailVerified: false,
				Status:        enum.UserStatusActive,
				Type:          enum.UserTypeCustomer,
			}

			err = userService.userRepository.CreateUser(tx, *user)
			if err != nil {
				return err
			}

			wallet := entity.Wallet{
				Balance: 0,
				UserID:  user.ID,
			}

			err = userService.walletRepository.CreateWallet(tx, wallet)
			if err != nil {
				return err
			}

			cart := entity.Cart{
				UserID: user.ID,
			}

			err = userService.cartRepository.CreateCart(tx, cart)
			if err != nil {
				return err
			}

			// userService.smsService.SendOTP(registerInfo.Phone, otp)
			return nil
		})
	}
	
	user, err = userService.userRepository.FindUserByPhone(userService.db, verifyAuthInfo.Phone)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}

	accessToken, refreshToken, err := userService.jwtService.GenerateToken(user.ID)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}

	return userdto.UserInfoResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Type:         user.Type.String(),
	}, nil
}

func (userService *UserService) FindUserByPhone(phone string) (*entity.User, error) {
	user, err := userService.userRepository.FindUserByPhone(userService.db, phone)
	if err != nil {
		return nil, err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return nil, notFoundError
	}
	return user, nil
}

func (userService *UserService) VerifyEmail(verifyInfo userdto.VerifyEmailRequest) error {
	var conflictErrors exception.ConflictErrors
	user, err := userService.GetUserByID(verifyInfo.UserID)
	if err != nil {
		return err
	}

	if user.EmailVerified {
		conflictErrors.Add(userService.constants.Field.Email, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	redisKey := userService.constants.RedisKey.GenerateOTPKey(verifyInfo.Email)
	err = userService.otpService.VerifyOTP(redisKey, verifyInfo.OTP)
	if err != nil {
		return err
	}
	user.EmailVerified = true
	err = userService.userRepository.UpdateUser(userService.db, *user)
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) AdminLogin(adminInfo userdto.AdminLoginRequest) (userdto.AdminInfoResponse, error) {
	user, err := userService.FindActiveUserByPhone(adminInfo.Phone)
	if err != nil {
		return userdto.AdminInfoResponse{}, err
	}
	
	if !user.IsAdmin {
		return userdto.AdminInfoResponse{},
		exception.NewAccessDeniedError("user is not admin", nil)
	}
	
	accessToken, refreshToken, err := userService.jwtService.GenerateToken(user.ID)
	if err != nil {
		return userdto.AdminInfoResponse{}, err
	}
	return userdto.AdminInfoResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		}, nil
}	
	
func (userService *UserService) GetDashboard() (userdto.DashboardResponse, error) {
	brandsCount, err := userService.brandRepository.GetBrandsCount(userService.db); 
	if err != nil {
		return userdto.DashboardResponse{}, err
	}
	
	categoriesCount, err := userService.categoryRepository.GetCategoriesCount(userService.db); 
	if err != nil {
		return userdto.DashboardResponse{}, err
	}
	
	productsCount, err := userService.productRepository.GetProductsCount(userService.db); 
	if err != nil {
		return userdto.DashboardResponse{}, err
	}
	
	return userdto.DashboardResponse{
		BrandsCount: brandsCount,
		CategoriesCount: categoriesCount,
		ProductsCount: productsCount,
		}, nil
}

func (userService *UserService) GetUserWalletBalance(userID uint) (userdto.UserWalletBalance, error) {
	wallet, err := userService.walletRepository.FindWalletByUserID(userService.db, userID)
	if err != nil {
		return userdto.UserWalletBalance{}, err
	}
	return userdto.UserWalletBalance{
		Balance: wallet.Balance,
	}, nil
}

func (userService *UserService) DepositWallet(balanceUpdateInfo userdto.UserBalanceUpdate) (userdto.UserWalletBalance, error) {
	var newBalance uint
	err := userService.db.WithTransaction(func(tx database.Database) error {
		wallet, err := userService.walletRepository.FindWalletByUserID(tx, balanceUpdateInfo.UserID)
		if err != nil {
			return err
		}

		transaction := entity.Transaction{
			Amount:   balanceUpdateInfo.Amount,
			WalletID: wallet.ID,
			Type:     enum.TransactionTypeDeposit,
		}
		if err = userService.transactionRepository.CreateTransaction(tx, transaction); err != nil {
			return err
		}
		
		newBalance, err = userService.walletRepository.DepositWallet(tx, balanceUpdateInfo.UserID, balanceUpdateInfo.Amount)
		if err != nil {
			return err
		}
		
		return nil
	})

	if err != nil {
		return userdto.UserWalletBalance{}, err
	}

	return userdto.UserWalletBalance{
		Balance: newBalance,
	}, nil
}

func (userService *UserService) WithdrawWallet(balanceUpdateInfo userdto.UserBalanceUpdate) (userdto.UserWalletBalance, error) {
	var newBalance uint
	var err error
	userService.db.WithTransaction(func(tx database.Database) error {
		wallet, err := userService.walletRepository.FindWalletByUserID(userService.db, balanceUpdateInfo.UserID)
		if err != nil {
			return err
		}

		transaction := entity.Transaction{
			Amount:   balanceUpdateInfo.Amount,
			WalletID: wallet.ID,
			Type:     enum.TransactionTypeWithDraw,
		}
		err = userService.transactionRepository.CreateTransaction(userService.db, transaction)
		if err != nil {
			return err
		}

		newBalance, err = userService.walletRepository.WithdrawWallet(userService.db, balanceUpdateInfo.UserID, balanceUpdateInfo.Amount)
		if err != nil {
			return err
		}

		return nil
	})

	return userdto.UserWalletBalance{
		Balance: newBalance,
	}, err
}