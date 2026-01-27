package service

import (
	"context"
	"regexp"
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
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	constants           *bootstrap.Constants
	otpService          usecase.OTPService
	jwtService          usecase.JWTService
	smsService          communication.SMSService
	emailService        communication.EmailService
	userRepository      postgres.UserRepository
	userCacheRepository redis.UserCacheRepository
	db                  database.Database
}

type UserServiceDeps struct {
	Constants           *bootstrap.Constants
	OTPService          usecase.OTPService
	JWTService          usecase.JWTService
	SMSService          communication.SMSService
	EmailService        communication.EmailService
	UserRepository      postgres.UserRepository
	UserCacheRepository redis.UserCacheRepository
	DB                  database.Database
}

func NewUserService(deps UserServiceDeps) *UserService {
	return &UserService{
		constants:           deps.Constants,
		otpService:          deps.OTPService,
		jwtService:          deps.JWTService,
		smsService:          deps.SMSService,
		emailService:        deps.EmailService,
		userRepository:      deps.UserRepository,
		userCacheRepository: deps.UserCacheRepository,
		db:                  deps.DB,
	}
}

func (userService *UserService) validatePasswordTests(errors *[]string, test string, password string, tag string) {
	matched, _ := regexp.MatchString(test, password)
	if !matched {
		*errors = append(*errors, tag)
	}
}

func (userService *UserService) passwordValidation(password string) error {
	var errors exception.ValidationErrors
	var errorTags []string

	userService.validatePasswordTests(&errorTags, ".{8,}", password, userService.constants.Tag.MinimumLength)
	userService.validatePasswordTests(&errorTags, "[a-z]", password, userService.constants.Tag.ContainsLowercase)
	userService.validatePasswordTests(&errorTags, "[A-Z]", password, userService.constants.Tag.ContainsUppercase)
	userService.validatePasswordTests(&errorTags, "[0-9]", password, userService.constants.Tag.ContainsNumber)
	userService.validatePasswordTests(&errorTags, "[^\\d\\w]", password, userService.constants.Tag.ContainsSpecialChar)

	for _, tag := range errorTags {
		errors.Add(userService.constants.Field.Password, tag)
	}
	if len(errorTags) > 0 {
		return errors
	}
	return nil
}

func (userService *UserService) validateDuplicateEmail(email string) error {
	var conflictErrors exception.ConflictErrors
	redisKey := userService.constants.RedisKey.GenerateOTPKey(email)
	data, err := userService.userCacheRepository.Get(context.Background(), redisKey)
	if err != nil {
		return err
	}
	if data != nil {
		conflictErrors.Add(userService.constants.Field.Email, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	user, err := userService.userRepository.FindUserByEmail(userService.db, email)
	if err != nil {
		return err
	}
	if user != nil && user.EmailVerified {
		conflictErrors.Add(userService.constants.Field.Email, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	return nil
}

func (userService *UserService) validateDuplicatePhone(phone string) error {
	var conflictErrors exception.ConflictErrors
	redisKey := userService.constants.RedisKey.GenerateOTPKey(phone)
	data, err := userService.userCacheRepository.Get(context.Background(), redisKey)
	if err != nil {
		return err
	}
	if data != nil {
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	user, err := userService.userRepository.FindUserByPhone(userService.db, phone)
	if err != nil {
		return err
	}
	if user != nil && user.PhoneVerified {
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.AlreadyRegistered)
		return conflictErrors
	}

	return nil
}

func (userService *UserService) enterNewEmail(firstName, lastName, email, emailSubject, templateFile string) error {
	err := userService.validateDuplicateEmail(email)
	if err != nil {
		return err
	}

	otp, expiryMinute, err := userService.otpService.GenerateOTP()
	if err != nil {
		return err
	}
	redisKey := userService.constants.RedisKey.GenerateOTPKey(email)
	err = userService.userCacheRepository.Set(context.Background(), redisKey, otp, time.Duration(expiryMinute)*time.Minute)
	if err != nil {
		return err
	}

	data := struct {
		FirstName    string
		LastName     string
		OTP          string
		ExpiryMinute int
		Year         int
	}{
		FirstName:    firstName,
		LastName:     lastName,
		OTP:          otp,
		ExpiryMinute: expiryMinute,
		Year:         time.Now().Year(),
	}
	if err := userService.emailService.SendEmail(email, emailSubject, templateFile, data); err != nil {
		return err
	}
	return nil
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
	if !user.PhoneVerified {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotVerified)
		return nil, conflictErrors
	}

	return user, nil
}

func (userService *UserService) GetUserCredential(userID uint) (userdto.CredentialResponse, error) {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return userdto.CredentialResponse{}, err
	}

	return userdto.CredentialResponse{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Phone:      user.Phone,
		Email:      user.Email,
		NationalID: user.NationalCode,
		Status:     user.Status.String(),
	}, nil
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
	err = userService.userRepository.UpdateUser(userService.db, user)
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
	err = userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) Register(registerInfo userdto.BasicRegisterRequest) error {
	err := userService.validateDuplicatePhone(registerInfo.Phone)
	if err != nil {
		return err
	}

	err = userService.passwordValidation(registerInfo.Password)
	if err != nil {
		return err
	}

	hashesPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(registerInfo.Password), 14)
	if err != nil {
		return err
	}

	err = userService.db.WithTransaction(func(tx database.Database) error {
		err = userService.userRepository.DeleteUserByPhone(tx, registerInfo.Phone)
		if err != nil {
			return err
		}

		user := &entity.User{
			FirstName:     registerInfo.FirstName,
			LastName:      registerInfo.LastName,
			Phone:         registerInfo.Phone,
			Password:      string(hashesPasswordBytes),
			PhoneVerified: false,
			EmailVerified: false,
			Status:        enum.UserStatusActive,
		}
		err = userService.userRepository.CreateUser(tx, user)
		if err != nil {
			return err
		}

		otp, expiryMinute, err := userService.otpService.GenerateOTP()
		if err != nil {
			return err
		}
		redisKey := userService.constants.RedisKey.GenerateOTPKey(registerInfo.Phone)
		err = userService.userCacheRepository.Set(context.Background(), redisKey, otp, time.Duration(expiryMinute)*time.Minute)
		if err != nil {
			return err
		}

		// userService.smsService.SendOTP(registerInfo.Phone, otp)
		return nil
	})

	return err
}

func (userService *UserService) VerifyPhone(verifyInfo userdto.VerifyPhoneRequest) error {
	user, err := userService.FindUserByPhone(verifyInfo.Phone)
	if err != nil {
		return err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return notFoundError
	}

	redisKey := userService.constants.RedisKey.GenerateOTPKey(verifyInfo.Phone)
	err = userService.otpService.VerifyOTP(redisKey, verifyInfo.OTP)
	if err != nil {
		return err
	}
	user.PhoneVerified = true
	err = userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) Login(loginInfo userdto.LoginRequest) (userdto.UserInfoResponse, error) {
	user, err := userService.FindActiveUserByPhone(loginInfo.Phone)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password))
	if err != nil {
		authError := exception.NewInvalidCredentialsError("phone and password not match", nil)
		return userdto.UserInfoResponse{}, authError
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
	}, nil
}

func (userService *UserService) ForgotPassword(forgotPasswordInfo userdto.ForgotPasswordRequest) error {
	_, err := userService.FindActiveUserByPhone(forgotPasswordInfo.Phone)
	if err != nil {
		return err
	}

	otp, expiryMinute, err := userService.otpService.GenerateOTP()
	if err != nil {
		return err
	}
	redisKey := userService.constants.RedisKey.GenerateOTPKey(forgotPasswordInfo.Phone)
	err = userService.userCacheRepository.Set(context.Background(), redisKey, otp, time.Duration(expiryMinute)*time.Minute)
	if err != nil {
		return err
	}
	// userService.smsService.SendOTP(registerInfo.Phone, otp)
	return nil
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

func (userService *UserService) VerifyOTP(verifyInfo userdto.VerifyPhoneRequest) (userdto.UserInfoResponse, error) {
	user, err := userService.FindUserByPhone(verifyInfo.Phone)
	if err != nil {
		return userdto.UserInfoResponse{}, err
	}

	redisKey := userService.constants.RedisKey.GenerateOTPKey(verifyInfo.Phone)
	err = userService.otpService.VerifyOTP(redisKey, verifyInfo.OTP)
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
	}, nil
}

func (userService *UserService) CompleteRegister(completeRegisterInfo userdto.CompleteRegisterRequest) error {
	user, err := userService.GetUserByID(completeRegisterInfo.UserID)
	if err != nil {
		return err
	}

	if completeRegisterInfo.Email != "" {
		userService.enterNewEmail(user.FirstName, user.LastName, completeRegisterInfo.Email, completeRegisterInfo.EmailSubject, completeRegisterInfo.TemplateFile)
	}
	user.Email = completeRegisterInfo.Email
	user.EmailVerified = false
	user.NationalCode = completeRegisterInfo.NationalCode

	err = userService.db.WithTransaction(func(tx database.Database) error {
		err = userService.userRepository.UpdateUser(tx, user)
		if err != nil {
			return err
		}
		return nil
	})
	return err
}

func (userService *UserService) VerifyEmail(verifyInfo userdto.VerifyEmailRequest) error {
	var conflictErrors exception.ConflictErrors
	user, err := userService.GetUserByID(verifyInfo.UserID)
	if err != nil {
		return err
	}

	if !user.PhoneVerified {
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotVerified)
		return conflictErrors
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
	err = userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) ResetPassword(resetPassInfo userdto.ResetPasswordRequest) error {
	user, err := userService.GetUserByID(resetPassInfo.UserID)
	if err != nil {
		return err
	}

	if !user.PhoneVerified {
		var conflictErrors exception.ConflictErrors
		conflictErrors.Add(userService.constants.Field.Phone, userService.constants.Tag.NotVerified)
		return conflictErrors
	}

	if err := userService.passwordValidation(resetPassInfo.Password); err != nil {
		return err
	}

	hashesPasswordBytes, err := bcrypt.GenerateFromPassword([]byte(resetPassInfo.Password), 14)
	if err != nil {
		return err
	}
	user.Password = string(hashesPasswordBytes)

	err = userService.userRepository.UpdateUser(userService.db, user)
	if err != nil {
		return err
	}
	return nil
}

func (userService *UserService) UpdateProfile(profileInfo userdto.UpdateProfileRequest) error {
	user, err := userService.GetUserByID(profileInfo.UserID)
	if err != nil {
		return err
	}

	if profileInfo.FirstName != nil {
		user.FirstName = *profileInfo.FirstName
	}

	if profileInfo.LastName != nil {
		user.LastName = *profileInfo.LastName
	}

	if profileInfo.Email != nil && user.Email != *profileInfo.Email {
		userService.enterNewEmail(user.FirstName, user.LastName, *profileInfo.Email, profileInfo.EmailSubject, profileInfo.TemplateFile)
		user.Email = *profileInfo.Email
		user.EmailVerified = false
	}

	if profileInfo.NationalCode != nil {
		user.NationalCode = *profileInfo.NationalCode
	}

	err = userService.db.WithTransaction(func(tx database.Database) error {
		if err := userService.userRepository.UpdateUser(tx, user); err != nil {
			return err
		}


		return nil
	})

	return err
}
