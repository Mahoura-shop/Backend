package mocks

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/stretchr/testify/mock"
)

type CurrencyRepositoryMock struct {
	mock.Mock
}

func NewCurrencyRepositoryMock() *CurrencyRepositoryMock {
	return &CurrencyRepositoryMock{}
}

func (m *CurrencyRepositoryMock) FindCurrencyByID(db database.Database, id uint) (*entity.Currency, error) {
	args := m.Called(db, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Currency), args.Error(1)
}

func (m *CurrencyRepositoryMock) FindCurrencyByCode(db database.Database, code string) (*entity.Currency, error) {
	args := m.Called(db, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Currency), args.Error(1)
}

func (m *CurrencyRepositoryMock) GetCurrencies(db database.Database) ([]*entity.Currency, error) {
	args := m.Called(db)
	return args.Get(0).([]*entity.Currency), args.Error(1)
}

func (m *CurrencyRepositoryMock) UpdateCurrency(db database.Database, currency entity.Currency) error {
	args := m.Called(db, currency)
	return args.Error(0)
}

func (m *CurrencyRepositoryMock) CreateCurrency(db database.Database, currency entity.Currency) error {
	args := m.Called(db, currency)
	return args.Error(0)
}
