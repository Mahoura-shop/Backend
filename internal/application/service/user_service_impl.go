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
	constants              *bootstrap.Constants
	otpService             usecase.OTPService
	jwtService             usecase.JWTService
	smsService             communication.SMSService
	emailService           communication.EmailService
	userRepository         postgres.UserRepository
	roleRepository         postgres.RoleRepository
	categoryRepository     postgres.CategoryRepository
	brandRepository        postgres.BrandRepository
	productRepository      postgres.ProductRepository
	productVisitRepository postgres.ProductVisitRepository
	walletRepository       postgres.WalletRepository
	cartRepository         postgres.CartRepository
	transactionRepository  postgres.TransactionRepository
	orderRepository        postgres.OrderRepository
	userCacheRepository    redis.UserCacheRepository
	db                     database.Database
	rbac                   *bootstrap.RBAC
}

type UserServiceDeps struct {
	Constants              *bootstrap.Constants
	OTPService             usecase.OTPService
	JWTService             usecase.JWTService
	SMSService             communication.SMSService
	EmailService           communication.EmailService
	UserRepository         postgres.UserRepository
	RoleRepository         postgres.RoleRepository
	CategoryRepository     postgres.CategoryRepository
	BrandRepository        postgres.BrandRepository
	ProductRepository      postgres.ProductRepository
	ProductVisitRepository postgres.ProductVisitRepository
	WalletRepository       postgres.WalletRepository
	CartRepository         postgres.CartRepository
	TransactionRepository  postgres.TransactionRepository
	OrderRepository        postgres.OrderRepository
	UserCacheRepository    redis.UserCacheRepository
	DB                     database.Database
	RBAC                   *bootstrap.RBAC
}

func NewUserService(deps UserServiceDeps) *UserService {
	return &UserService{
		constants:              deps.Constants,
		otpService:             deps.OTPService,
		jwtService:             deps.JWTService,
		smsService:             deps.SMSService,
		emailService:           deps.EmailService,
		userRepository:         deps.UserRepository,
		roleRepository:         deps.RoleRepository,
		categoryRepository:     deps.CategoryRepository,
		brandRepository:        deps.BrandRepository,
		productRepository:      deps.ProductRepository,
		productVisitRepository: deps.ProductVisitRepository,
		walletRepository:       deps.WalletRepository,
		cartRepository:         deps.CartRepository,
		transactionRepository:  deps.TransactionRepository,
		orderRepository:        deps.OrderRepository,
		userCacheRepository:    deps.UserCacheRepository,
		db:                     deps.DB,
		rbac:                   deps.RBAC,
	}
}


func (userService *UserService) ParseUser(user entity.User) userdto.UserCredential {
	roleName := ""
	if user.Role != nil {
		roleName = user.Role.Name
	}
	return userdto.UserCredential{
		ID:         user.ID,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Phone:      user.Phone,
		Email:      user.Email,
		ProfilePic: user.ProfilePicPath,
		Status:     user.Status.String(),
		Type:       user.Type.String(),
		IsAdmin:    user.IsAdmin,
		RoleID:     user.RoleID,
		RoleName:   roleName,
		CreatedAt:  user.CreatedAt.Format(time.RFC3339),
	}
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
	if err := userService.smsService.SendOTP(authInfo.Phone, otp); err != nil {
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

	permissions := []string{}
	if user.IsAdmin {
		if !userService.rbac.UseRBAC || user.RoleID == nil {
			allPerms, permsErr := userService.roleRepository.GetPermissions(userService.db)
			if permsErr == nil {
				for _, p := range allPerms {
					permissions = append(permissions, p.Name)
				}
			}
		} else {
			role, roleErr := userService.roleRepository.FindRoleByID(userService.db, *user.RoleID)
			if roleErr == nil && role != nil {
				for _, p := range role.Permissions {
					permissions = append(permissions, p.Name)
				}
			}
		}
	}

	return userdto.UserInfoResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Type:         user.Type.String(),
		IsAdmin:      user.IsAdmin,
		Permissions:  permissions,
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



func (userService *UserService) GetOrdersChart(period string) ([]userdto.OrderByDay, error) {
	days := 7
	switch period {
	case "month":
		days = 30
	case "year":
		days = 365
	}
	orders, err := userService.orderRepository.GetOrdersPerDay(userService.db, days)
	if err != nil {
		return nil, err
	}
	countByDate := make(map[string]uint, len(orders))
	for _, o := range orders {
		key := o.Date
		if len(key) > 10 {
			key = key[:10]
		}
		countByDate[key] = o.Count
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	result := make([]userdto.OrderByDay, days)
	for i := range result {
		d := today.AddDate(0, 0, -(days-1-i))
		dateStr := d.Format("2006-01-02")
		result[i] = userdto.OrderByDay{Date: dateStr, Count: countByDate[dateStr]}
	}
	return result, nil
}

func (userService *UserService) GetSalesChart(period string) ([]userdto.RevenueByDay, error) {
	days := 7
	switch period {
	case "month":
		days = 30
	case "year":
		days = 365
	}
	rows, err := userService.orderRepository.GetRevenuePerDay(userService.db, days)
	if err != nil {
		return nil, err
	}
	revenueByDate := make(map[string]uint, len(rows))
	for _, r := range rows {
		key := r.Date
		if len(key) > 10 {
			key = key[:10]
		}
		revenueByDate[key] = r.Revenue
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	result := make([]userdto.RevenueByDay, days)
	for i := range result {
		d := today.AddDate(0, 0, -(days-1-i))
		dateStr := d.Format("2006-01-02")
		result[i] = userdto.RevenueByDay{Date: dateStr, Revenue: revenueByDate[dateStr]}
	}
	return result, nil
}

func (userService *UserService) GetVisitsChart(period string) ([]userdto.VisitsByDay, error) {
	days := 7
	switch period {
	case "month":
		days = 30
	case "year":
		days = 365
	}
	visits, err := userService.productVisitRepository.GetAllVisitsPerDay(userService.db, days)
	if err != nil {
		return nil, err
	}
	countByDate := make(map[string]uint, len(visits))
	for _, v := range visits {
		key := v.Date
		if len(key) > 10 {
			key = key[:10]
		}
		countByDate[key] = v.Count
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	result := make([]userdto.VisitsByDay, days)
	for i := range result {
		d := today.AddDate(0, 0, -(days-1-i))
		dateStr := d.Format("2006-01-02")
		result[i] = userdto.VisitsByDay{Date: dateStr, Count: countByDate[dateStr]}
	}
	return result, nil
}

func (userService *UserService) GetProductVisitsChart(productID uint, period string) ([]userdto.VisitsByDay, error) {
	days := 7
	switch period {
	case "month":
		days = 30
	case "year":
		days = 365
	}
	visits, err := userService.productVisitRepository.GetVisitsPerDay(userService.db, productID, days)
	if err != nil {
		return nil, err
	}
	countByDate := make(map[string]uint, len(visits))
	for _, v := range visits {
		key := v.Date
		if len(key) > 10 {
			key = key[:10]
		}
		countByDate[key] = v.Count
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	result := make([]userdto.VisitsByDay, days)
	for i := range result {
		d := today.AddDate(0, 0, -(days-1-i))
		dateStr := d.Format("2006-01-02")
		result[i] = userdto.VisitsByDay{Date: dateStr, Count: countByDate[dateStr]}
	}
	return result, nil
}

func (userService *UserService) GetProductOrdersChart(productID uint, period string) ([]userdto.OrderByDay, error) {
	days := 7
	switch period {
	case "month":
		days = 30
	case "year":
		days = 365
	}
	orders, err := userService.orderRepository.GetProductOrdersPerDay(userService.db, productID, days)
	if err != nil {
		return nil, err
	}
	countByDate := make(map[string]uint, len(orders))
	for _, o := range orders {
		key := o.Date
		if len(key) > 10 {
			key = key[:10]
		}
		countByDate[key] = o.Count
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	result := make([]userdto.OrderByDay, days)
	for i := range result {
		d := today.AddDate(0, 0, -(days-1-i))
		dateStr := d.Format("2006-01-02")
		result[i] = userdto.OrderByDay{Date: dateStr, Count: countByDate[dateStr]}
	}
	return result, nil
}

func (userService *UserService) GetCategoryVisitsChart(categoryID uint, period string) ([]userdto.VisitsByDay, error) {
	days := 7
	switch period {
	case "month":
		days = 30
	case "year":
		days = 365
	}
	visits, err := userService.productVisitRepository.GetCategoryVisitsPerDay(userService.db, categoryID, days)
	if err != nil {
		return nil, err
	}
	countByDate := make(map[string]uint, len(visits))
	for _, v := range visits {
		key := v.Date
		if len(key) > 10 {
			key = key[:10]
		}
		countByDate[key] = v.Count
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	result := make([]userdto.VisitsByDay, days)
	for i := range result {
		d := today.AddDate(0, 0, -(days-1-i))
		dateStr := d.Format("2006-01-02")
		result[i] = userdto.VisitsByDay{Date: dateStr, Count: countByDate[dateStr]}
	}
	return result, nil
}

func (userService *UserService) GetCategoryOrdersChart(categoryID uint, period string) ([]userdto.OrderByDay, error) {
	days := 7
	switch period {
	case "month":
		days = 30
	case "year":
		days = 365
	}
	orders, err := userService.orderRepository.GetCategoryOrdersPerDay(userService.db, categoryID, days)
	if err != nil {
		return nil, err
	}
	countByDate := make(map[string]uint, len(orders))
	for _, o := range orders {
		key := o.Date
		if len(key) > 10 {
			key = key[:10]
		}
		countByDate[key] = o.Count
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	result := make([]userdto.OrderByDay, days)
	for i := range result {
		d := today.AddDate(0, 0, -(days-1-i))
		dateStr := d.Format("2006-01-02")
		result[i] = userdto.OrderByDay{Date: dateStr, Count: countByDate[dateStr]}
	}
	return result, nil
}

func (userService *UserService) GetBrandVisitsChart(brandID uint, period string) ([]userdto.VisitsByDay, error) {
	days := 7
	switch period {
	case "month":
		days = 30
	case "year":
		days = 365
	}
	visits, err := userService.productVisitRepository.GetBrandVisitsPerDay(userService.db, brandID, days)
	if err != nil {
		return nil, err
	}
	countByDate := make(map[string]uint, len(visits))
	for _, v := range visits {
		key := v.Date
		if len(key) > 10 {
			key = key[:10]
		}
		countByDate[key] = v.Count
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	result := make([]userdto.VisitsByDay, days)
	for i := range result {
		d := today.AddDate(0, 0, -(days-1-i))
		dateStr := d.Format("2006-01-02")
		result[i] = userdto.VisitsByDay{Date: dateStr, Count: countByDate[dateStr]}
	}
	return result, nil
}

func (userService *UserService) GetBrandOrdersChart(brandID uint, period string) ([]userdto.OrderByDay, error) {
	days := 7
	switch period {
	case "month":
		days = 30
	case "year":
		days = 365
	}
	orders, err := userService.orderRepository.GetBrandOrdersPerDay(userService.db, brandID, days)
	if err != nil {
		return nil, err
	}
	countByDate := make(map[string]uint, len(orders))
	for _, o := range orders {
		key := o.Date
		if len(key) > 10 {
			key = key[:10]
		}
		countByDate[key] = o.Count
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	result := make([]userdto.OrderByDay, days)
	for i := range result {
		d := today.AddDate(0, 0, -(days-1-i))
		dateStr := d.Format("2006-01-02")
		result[i] = userdto.OrderByDay{Date: dateStr, Count: countByDate[dateStr]}
	}
	return result, nil
}

func (userService *UserService) GetPublicStats() (userdto.PublicStatsResponse, error) {
	productsCount, err := userService.productRepository.GetProductsCount(userService.db)
	if err != nil {
		return userdto.PublicStatsResponse{}, err
	}

	brandsCount, err := userService.brandRepository.GetBrandsCount(userService.db)
	if err != nil {
		return userdto.PublicStatsResponse{}, err
	}

	usersCount, err := userService.userRepository.GetUsersCount(userService.db)
	if err != nil {
		return userdto.PublicStatsResponse{}, err
	}

	return userdto.PublicStatsResponse{
		ProductsCount: productsCount,
		BrandsCount:   brandsCount,
		UsersCount:    usersCount,
	}, nil
}

func (userService *UserService) GetDashboard() (userdto.DashboardResponse, error) {
	brandsCount, err := userService.brandRepository.GetBrandsCount(userService.db)
	if err != nil {
		return userdto.DashboardResponse{}, err
	}

	categoriesCount, err := userService.categoryRepository.GetCategoriesCount(userService.db)
	if err != nil {
		return userdto.DashboardResponse{}, err
	}

	productsCount, err := userService.productRepository.GetProductsCount(userService.db)
	if err != nil {
		return userdto.DashboardResponse{}, err
	}

	revenuePerTier, err := userService.orderRepository.GetRevenuePerTier(userService.db)
	if err != nil {
		return userdto.DashboardResponse{}, err
	}

	ordersPerDay, err := userService.orderRepository.GetOrdersPerDay(userService.db, 7)
	if err != nil {
		return userdto.DashboardResponse{}, err
	}

	lowStockProducts, err := userService.productRepository.GetLowStockProducts(userService.db, 5)
	if err != nil {
		return userdto.DashboardResponse{}, err
	}

	topProducts, err := userService.productRepository.GetTopSoldProducts(userService.db, 5)
	if err != nil {
		return userdto.DashboardResponse{}, err
	}

	revenuePerTierDTO := make([]userdto.RevenueByTier, len(revenuePerTier))
	for i, r := range revenuePerTier {
		revenuePerTierDTO[i] = userdto.RevenueByTier{
			Tier:    r.Tier,
			Revenue: r.Revenue,
		}
	}

	ordersPerDayDTO := make([]userdto.OrderByDay, len(ordersPerDay))
	for i, o := range ordersPerDay {
		ordersPerDayDTO[i] = userdto.OrderByDay{
			Date:  o.Date,
			Count: o.Count,
		}
	}

	lowStockDTO := make([]userdto.LowStockProduct, len(lowStockProducts))
	for i, p := range lowStockProducts {
		lowStockDTO[i] = userdto.LowStockProduct{
			ID:       p.ID,
			Name:     p.Name,
			Quantity: p.Quantity,
			MinOrder: p.MinOrder,
		}
	}

	topProductsDTO := make([]userdto.TopProduct, len(topProducts))
	for i, p := range topProducts {
		topProductsDTO[i] = userdto.TopProduct{
			ID:       p.ID,
			Name:     p.Name,
			Quantity: p.Quantity,
			Revenue:  p.Revenue,
		}
	}

	return userdto.DashboardResponse{
		ProductsCount:    productsCount,
		CategoriesCount:  categoriesCount,
		BrandsCount:      brandsCount,
		RevenuePerTier:   revenuePerTierDTO,
		OrdersPerDay:     ordersPerDayDTO,
		LowStockProducts: lowStockDTO,
		TopProducts:      topProductsDTO,
	}, nil
}

func (userService *UserService) GetProvinceStats() ([]userdto.ProvinceStatDTO, error) {
	rows, err := userService.orderRepository.GetProvinceStats(userService.db)
	if err != nil {
		return nil, err
	}
	result := make([]userdto.ProvinceStatDTO, len(rows))
	for i, r := range rows {
		result[i] = userdto.ProvinceStatDTO{
			Province:   r.Province,
			OrderCount: r.OrderCount,
			Revenue:    r.Revenue,
		}
	}
	return result, nil
}

func (userService *UserService) GetUserWalletBalance(userID uint) (userdto.UserWalletBalance, error) {
	wallet, err := userService.walletRepository.FindWalletByUserID(userService.db, userID)
	if err != nil {
		return userdto.UserWalletBalance{}, err
	}
	if wallet == nil {
		return userdto.UserWalletBalance{Balance: 0}, nil
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
		if wallet == nil {
			if err = userService.walletRepository.CreateWallet(tx, entity.Wallet{UserID: balanceUpdateInfo.UserID}); err != nil {
				return err
			}
			wallet, err = userService.walletRepository.FindWalletByUserID(tx, balanceUpdateInfo.UserID)
			if err != nil {
				return err
			}
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
		wallet, err := userService.walletRepository.FindWalletByUserID(tx, balanceUpdateInfo.UserID)
		if err != nil {
			return err
		}
		if wallet == nil {
			return nil
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

func (userService *UserService) UpdateProfile(req userdto.UpdateProfileRequest) (userdto.UserCredential, error) {
	user, err := userService.userRepository.FindUserByID(userService.db, req.UserID)
	if err != nil {
		return userdto.UserCredential{}, err
	}
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	err = userService.userRepository.UpdateUser(userService.db, *user)
	if err != nil {
		return userdto.UserCredential{}, err
	}
	return userService.ParseUser(*user), nil
}

func (userService *UserService) GetWalletHistory(userID uint) ([]userdto.TransactionDTO, error) {
	wallet, err := userService.walletRepository.FindWalletByUserID(userService.db, userID)
	if err != nil {
		return nil, err
	}
	if wallet == nil {
		return []userdto.TransactionDTO{}, nil
	}

	transactions, err := userService.transactionRepository.FindTransactionsByWalletID(userService.db, wallet.ID)
	if err != nil {
		return nil, err
	}

	dtos := make([]userdto.TransactionDTO, len(transactions))
	for i, t := range transactions {
		dtos[i] = userdto.TransactionDTO{
			ID:        t.ID,
			Amount:    t.Amount,
			Type:      uint(t.Type),
			CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return dtos, nil
}

func (userService *UserService) GetAdminUserWallet(userID uint) (userdto.AdminUserWalletResponse, error) {
	wallet, err := userService.walletRepository.FindWalletByUserID(userService.db, userID)
	if err != nil {
		return userdto.AdminUserWalletResponse{}, err
	}
	if wallet == nil {
		return userdto.AdminUserWalletResponse{Balance: 0, Transactions: []userdto.TransactionDTO{}}, nil
	}

	transactions, err := userService.transactionRepository.FindTransactionsByWalletID(userService.db, wallet.ID)
	if err != nil {
		return userdto.AdminUserWalletResponse{}, err
	}

	dtos := make([]userdto.TransactionDTO, len(transactions))
	for i, t := range transactions {
		dtos[i] = userdto.TransactionDTO{
			ID:        t.ID,
			Amount:    t.Amount,
			Type:      uint(t.Type),
			CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z"),
		}
	}

	return userdto.AdminUserWalletResponse{
		Balance:      wallet.Balance,
		Transactions: dtos,
	}, nil
}

func (userService *UserService) GetSubAdmins() ([]userdto.SubAdminCredential, error) {
	users, err := userService.userRepository.FindAdmins(userService.db)
	if err != nil {
		return nil, err
	}
	result := make([]userdto.SubAdminCredential, len(users))
	for i, u := range users {
		cred := userdto.SubAdminCredential{
			ID:        u.ID,
			FirstName: u.FirstName,
			LastName:  u.LastName,
			Phone:     u.Phone,
			Email:     u.Email,
			RoleID:    u.RoleID,
		}
		if u.Role != nil {
			cred.RoleName = u.Role.Name
		}
		result[i] = cred
	}
	return result, nil
}

func (userService *UserService) CreateSubAdmin(phone string, roleID uint) error {
	user, err := userService.userRepository.FindUserByPhone(userService.db, phone)
	if err != nil {
		return err
	}
	if user == nil {
		notFoundError := exception.NotFoundError{Item: userService.constants.Field.User}
		return notFoundError
	}
	user.IsAdmin = true
	user.RoleID = &roleID
	return userService.userRepository.UpdateUser(userService.db, *user)
}

func (userService *UserService) AssignSubAdminRole(userID uint, roleID uint) error {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return err
	}
	user.RoleID = &roleID
	return userService.userRepository.UpdateUser(userService.db, *user)
}

func (userService *UserService) RevokeSubAdmin(userID uint) error {
	user, err := userService.GetUserByID(userID)
	if err != nil {
		return err
	}
	user.IsAdmin = false
	user.RoleID = nil
	return userService.userRepository.UpdateUser(userService.db, *user)
}