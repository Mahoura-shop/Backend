package service

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	currencydto "github.com/Mahoura-shop/Backend/internal/application/dto/currency"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	"github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/domain/s3"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type CurrencyService struct {
	constants          *bootstrap.Constants
	currencyRepository postgres.CurrencyRepository
	s3Storage          s3.S3Storage
	db                 database.Database
}

type CurrencyServiceDeps struct {
	Constants          *bootstrap.Constants
	CurrencyRepository postgres.CurrencyRepository
	S3Storage          s3.S3Storage
	DB                 database.Database
}

func NewCurrencyService(deps CurrencyServiceDeps) *CurrencyService {
	return &CurrencyService{
		constants:          deps.Constants,
		currencyRepository: deps.CurrencyRepository,
		s3Storage:          deps.S3Storage,
		db:                 deps.DB,
	}
}

func (currencyService *CurrencyService) ParseCurrency(currency entity.Currency) (currencydto.CurrencyCredential, error) {
	response := currencydto.CurrencyCredential{
		ID:          currency.ID,
		Name:        currency.Name,
		Code:        currency.Code,
		ConvertRate: currency.ConvertRate,
	}
	return response, nil
}

func (currencyService *CurrencyService) FindCurrencyByCode(code string) (*currencydto.CurrencyCredential, error) {
	currency, err := currencyService.currencyRepository.FindCurrencyByCode(currencyService.db, code)
	if err != nil {
		return nil, err
	}
	if currency == nil {
		notFoundError := exception.NotFoundError{Item: currencyService.constants.Field.Currency}
		return nil, notFoundError
	}

	parsedCurrency, err := currencyService.ParseCurrency(*currency)
	if err != nil {
		return nil, err
	}
	return &parsedCurrency, nil
}

func (currencyService *CurrencyService) FindCurrencyByID(currencyID uint) (*currencydto.CurrencyCredential, error) {
	currency, err := currencyService.currencyRepository.FindCurrencyByID(currencyService.db, currencyID)
	if err != nil {
		return nil, err
	}
	if currency == nil {
		notFoundError := exception.NotFoundError{Item: currencyService.constants.Field.Currency}
		return nil, notFoundError
	}

	parsedCurrency, err := currencyService.ParseCurrency(*currency)
	if err != nil {
		return nil, err
	}
	return &parsedCurrency, nil
}

func (currencyService *CurrencyService) GetCurrencies() ([]currencydto.CurrencyCredential, error) {
	currencies, err := currencyService.currencyRepository.GetCurrencies(currencyService.db)
	if err != nil {
		return nil, err
	}
	var responses []currencydto.CurrencyCredential
	for _, currency := range currencies {
		response, err := currencyService.ParseCurrency(*currency)
		if err != nil {
			return nil, err
		}
		responses = append(responses, response)
	}
	return responses, nil
}

func (currencyService *CurrencyService) applyCurrencyUpdates(currency *entity.Currency, name *string, code *string, convertRate *uint) {
	if name != nil {
		currency.Name = *name
	}

	if code != nil {
		currency.Code = *code
	}

	if convertRate != nil {
		currency.ConvertRate = *convertRate
	}
}

func (currencyService *CurrencyService) UpdateCurrency(currencyInfo currencydto.UpdateCurrencyRequest) error {
	currency, err := currencyService.currencyRepository.FindCurrencyByID(currencyService.db, currencyInfo.ID)
	if err != nil {
		return err
	}
	if currency == nil {
		return exception.NotFoundError{Item: currencyService.constants.Field.Currency}
	}

	currencyService.applyCurrencyUpdates(currency, currencyInfo.Name, currencyInfo.Code, currencyInfo.ConvertRate)

	err = currencyService.db.WithTransaction(func(tx database.Database) error {
		if err := currencyService.currencyRepository.UpdateCurrency(tx, currency); err != nil {
			return err
		}
		return nil
	})

	return err
}