package mocks

import (
	userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *UserServiceMock {
	return &UserServiceMock{}
}

func (s *UserServiceMock) IsUserActive(userID uint) error {
	args := s.Called(userID)
	return args.Error(0)
}

func (s *UserServiceMock) GetUserByID(userID uint) (*entity.User, error) {
	args := s.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (s *UserServiceMock) GetUserCredential(userID uint) (userdto.UserCredential, error) {
	args := s.Called(userID)
	return args.Get(0).(userdto.UserCredential), args.Error(1)
}

func (s *UserServiceMock) BanUser(userID uint) error {
	args := s.Called(userID)
	return args.Error(0)
}

func (s *UserServiceMock) UnbanUser(userID uint) error {
	args := s.Called(userID)
	return args.Error(0)
}

func (s *UserServiceMock) VerifyEmail(verifyOTPInfo userdto.VerifyEmailRequest) error {
	args := s.Called(verifyOTPInfo)
	return args.Error(0)
}

func (s *UserServiceMock) FindActiveUserByPhone(phone string) (*entity.User, error) {
	args := s.Called(phone)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (s *UserServiceMock) ParseUser(user entity.User) userdto.UserCredential {
	args := s.Called(user)
	return args.Get(0).(userdto.UserCredential)
}

func (s *UserServiceMock) GetUsers() ([]userdto.UserCredential, error) {
	args := s.Called()
	return args.Get(0).([]userdto.UserCredential), args.Error(1)
}

func (s *UserServiceMock) Auth(req userdto.AuthRequest) error {
	args := s.Called(req)
	return args.Error(0)
}

func (s *UserServiceMock) VerifyAuth(req userdto.VerifyAuthRequest) (userdto.UserInfoResponse, error) {
	args := s.Called(req)
	return args.Get(0).(userdto.UserInfoResponse), args.Error(1)
}

func (s *UserServiceMock) FindUserByPhone(phone string) (*entity.User, error) {
	args := s.Called(phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}


func (s *UserServiceMock) GetPublicStats() (userdto.PublicStatsResponse, error) {
	args := s.Called()
	return args.Get(0).(userdto.PublicStatsResponse), args.Error(1)
}

func (s *UserServiceMock) GetDashboard() (userdto.DashboardResponse, error) {
	args := s.Called()
	return args.Get(0).(userdto.DashboardResponse), args.Error(1)
}

func (s *UserServiceMock) GetOrdersChart(period string) ([]userdto.OrderByDay, error) {
	args := s.Called(period)
	return args.Get(0).([]userdto.OrderByDay), args.Error(1)
}

func (s *UserServiceMock) GetSalesChart(period string) ([]userdto.RevenueByDay, error) {
	args := s.Called(period)
	return args.Get(0).([]userdto.RevenueByDay), args.Error(1)
}

func (s *UserServiceMock) GetProvinceStats() ([]userdto.ProvinceStatDTO, error) {
	args := s.Called()
	return args.Get(0).([]userdto.ProvinceStatDTO), args.Error(1)
}

func (s *UserServiceMock) GetUserWalletBalance(userID uint) (userdto.UserWalletBalance, error) {
	args := s.Called(userID)
	return args.Get(0).(userdto.UserWalletBalance), args.Error(1)
}

func (s *UserServiceMock) GetWalletHistory(userID uint) ([]userdto.TransactionDTO, error) {
	args := s.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdto.TransactionDTO), args.Error(1)
}

func (s *UserServiceMock) DepositWallet(req userdto.UserBalanceUpdate) (userdto.UserWalletBalance, error) {
	args := s.Called(req)
	return args.Get(0).(userdto.UserWalletBalance), args.Error(1)
}

func (s *UserServiceMock) WithdrawWallet(req userdto.UserBalanceUpdate) (userdto.UserWalletBalance, error) {
	args := s.Called(req)
	return args.Get(0).(userdto.UserWalletBalance), args.Error(1)
}

func (s *UserServiceMock) GetAdminUserWallet(userID uint) (userdto.AdminUserWalletResponse, error) {
	args := s.Called(userID)
	return args.Get(0).(userdto.AdminUserWalletResponse), args.Error(1)
}

func (s *UserServiceMock) UpdateProfile(req userdto.UpdateProfileRequest) (userdto.UserCredential, error) {
	args := s.Called(req)
	return args.Get(0).(userdto.UserCredential), args.Error(1)
}

func (s *UserServiceMock) GetSubAdmins() ([]userdto.SubAdminCredential, error) {
	args := s.Called()
	return args.Get(0).([]userdto.SubAdminCredential), args.Error(1)
}

func (s *UserServiceMock) CreateSubAdmin(phone string, roleID uint) error {
	args := s.Called(phone, roleID)
	return args.Error(0)
}

func (s *UserServiceMock) AssignSubAdminRole(userID uint, roleID uint) error {
	args := s.Called(userID, roleID)
	return args.Error(0)
}

func (s *UserServiceMock) RevokeSubAdmin(userID uint) error {
	args := s.Called(userID)
	return args.Error(0)
}

func (s *UserServiceMock) GetVisitsChart(period string) ([]userdto.VisitsByDay, error) {
	args := s.Called(period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdto.VisitsByDay), args.Error(1)
}

func (s *UserServiceMock) GetProductVisitsChart(productID uint, period string) ([]userdto.VisitsByDay, error) {
	args := s.Called(productID, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdto.VisitsByDay), args.Error(1)
}

func (s *UserServiceMock) GetProductOrdersChart(productID uint, period string) ([]userdto.OrderByDay, error) {
	args := s.Called(productID, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdto.OrderByDay), args.Error(1)
}

func (s *UserServiceMock) GetCategoryVisitsChart(categoryID uint, period string) ([]userdto.VisitsByDay, error) {
	args := s.Called(categoryID, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdto.VisitsByDay), args.Error(1)
}

func (s *UserServiceMock) GetCategoryOrdersChart(categoryID uint, period string) ([]userdto.OrderByDay, error) {
	args := s.Called(categoryID, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdto.OrderByDay), args.Error(1)
}

func (s *UserServiceMock) GetBrandVisitsChart(brandID uint, period string) ([]userdto.VisitsByDay, error) {
	args := s.Called(brandID, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdto.VisitsByDay), args.Error(1)
}

func (s *UserServiceMock) GetBrandOrdersChart(brandID uint, period string) ([]userdto.OrderByDay, error) {
	args := s.Called(brandID, period)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]userdto.OrderByDay), args.Error(1)
}