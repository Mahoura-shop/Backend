package mocks

import (
	currencydto "github.com/Mahoura-shop/Backend/internal/application/dto/currency"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/stretchr/testify/mock"
)

type CurrencyServiceMock struct {
	mock.Mock
}

func NewCurrencyServiceMock() *CurrencyServiceMock {
	return &CurrencyServiceMock{}
}

func (m *CurrencyServiceMock) ParseCurrency(c entity.Currency) (currencydto.CurrencyCredential, error) {
	args := m.Called(c)
	return args.Get(0).(currencydto.CurrencyCredential), args.Error(1)
}

func (m *CurrencyServiceMock) FindCurrencyByCode(code string) (*currencydto.CurrencyCredential, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*currencydto.CurrencyCredential), args.Error(1)
}

func (m *CurrencyServiceMock) FindCurrencyByID(id uint) (*currencydto.CurrencyCredential, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*currencydto.CurrencyCredential), args.Error(1)
}

func (m *CurrencyServiceMock) GetCurrencies() ([]currencydto.CurrencyCredential, error) {
	args := m.Called()
	return args.Get(0).([]currencydto.CurrencyCredential), args.Error(1)
}

func (m *CurrencyServiceMock) UpdateCurrency(req currencydto.UpdateCurrencyRequest) error {
	args := m.Called(req)
	return args.Error(0)
}
