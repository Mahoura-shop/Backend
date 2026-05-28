package service

import (
	"testing"

	walletdto "github.com/Mahoura-shop/Backend/internal/application/dto/wallet"
	userdto "github.com/Mahoura-shop/Backend/internal/application/dto/user"
	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/suite"
)

type WalletServiceTestSuite struct {
	suite.Suite
	walletRepo  *mocks.WalletRepositoryMock
	userService *mocks.UserServiceMock
	db          *mocks.DatabaseMock
	service     *WalletService
}

func (s *WalletServiceTestSuite) SetupTest() {
	s.walletRepo = mocks.NewWalletRepositoryMock()
	s.userService = mocks.NewUserServiceMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewWalletService(WalletServiceDeps{
		Constants:        bootstrap.NewConstants(),
		WalletRepository: s.walletRepo,
		UserService:      s.userService,
		DB:               s.db,
	})
}

func (s *WalletServiceTestSuite) TestParseWallet_MapsBalanceAndID() {
	user := entity.User{
		Model: dbmodel.Model{ID: 5},
		Phone: "+989123456789",
	}
	wallet := entity.Wallet{
		Model:   dbmodel.Model{ID: 10},
		Balance: 500_000,
		User:    user,
	}
	parsedUser := userdto.UserCredential{ID: 5, Phone: "+989123456789"}
	s.userService.On("ParseUser", user).Return(parsedUser).Once()

	result := s.service.ParseWallet(wallet)

	s.Equal(uint(10), result.ID)
	s.Equal(uint(500_000), result.Balance)
	s.Equal(parsedUser, result.User)
	s.userService.AssertExpectations(s.T())
}

func (s *WalletServiceTestSuite) TestParseWallet_ZeroBalance() {
	wallet := entity.Wallet{
		Model:   dbmodel.Model{ID: 1},
		Balance: 0,
	}
	s.userService.On("ParseUser", wallet.User).Return(userdto.UserCredential{}).Once()

	result := s.service.ParseWallet(wallet)

	s.Equal(uint(0), result.Balance)
	s.userService.AssertExpectations(s.T())
}

func (s *WalletServiceTestSuite) TestParseWallet_ReturnsCorrectType() {
	wallet := entity.Wallet{}
	s.userService.On("ParseUser", wallet.User).Return(userdto.UserCredential{}).Once()

	result := s.service.ParseWallet(wallet)

	_, ok := interface{}(result).(walletdto.WalletCredential)
	s.True(ok)
	s.userService.AssertExpectations(s.T())
}

func TestWalletServiceSuite(t *testing.T) {
	suite.Run(t, new(WalletServiceTestSuite))
}
