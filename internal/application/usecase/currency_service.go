package usecase

import (
	currencydto "github.com/Mahoura-shop/Backend/internal/application/dto/currency"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
)

type CurrencyService interface {
	ParseCurrency(entity.Currency) (currencydto.CurrencyCredential, error)
	FindCurrencyByCode(string) (*currencydto.CurrencyCredential, error) 
	FindCurrencyByID(uint) (*currencydto.CurrencyCredential, error) 
	GetCurrencies() ([]currencydto.CurrencyCredential, error)
	UpdateCurrency(currencydto.UpdateCurrencyRequest) error
}
