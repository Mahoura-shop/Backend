package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type WalletRepositoryMock struct {
	mock.Mock
}

func NewWalletRepositoryMock() *WalletRepositoryMock {
	return &WalletRepositoryMock{}
}

func (m *WalletRepositoryMock) FindWalletByID(db database.Database, id uint) (*entity.Wallet, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Wallet), args.Error(1)
}

func (m *WalletRepositoryMock) FindWalletByUserID(db database.Database, userID uint) (*entity.Wallet, error) {
	args := m.Called(db, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Wallet), args.Error(1)
}

func (m *WalletRepositoryMock) FindWalletByPhone(db database.Database, phone string) (*entity.Wallet, error) {
	args := m.Called(db, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Wallet), args.Error(1)
}

func (m *WalletRepositoryMock) CreateWallet(db database.Database, wallet entity.Wallet) error {
	args := m.Called(db, wallet)
	return args.Error(0)
}

func (m *WalletRepositoryMock) DeleteWalletByPhone(db database.Database, phone string) error {
	args := m.Called(db, phone)
	return args.Error(0)
}

func (m *WalletRepositoryMock) UpdateWallet(db database.Database, wallet entity.Wallet) error {
	args := m.Called(db, wallet)
	return args.Error(0)
}

func (m *WalletRepositoryMock) DepositWallet(db database.Database, userID uint, amount uint) (uint, error) {
	args := m.Called(db, userID, amount)
	return uint(args.Int(0)), args.Error(1)
}

func (m *WalletRepositoryMock) WithdrawWallet(db database.Database, userID uint, amount uint) (uint, error) {
	args := m.Called(db, userID, amount)
	return uint(args.Int(0)), args.Error(1)
}
