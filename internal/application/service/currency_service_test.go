package service

import (
	"errors"
	"testing"

	"github.com/Mahoura-shop/Backend/bootstrap"
	currencydto "github.com/Mahoura-shop/Backend/internal/application/dto/currency"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/exception"
	dbmodel "github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type CurrencyServiceTestSuite struct {
	suite.Suite
	constants          *bootstrap.Constants
	currencyRepository *mocks.CurrencyRepositoryMock
	db                 *mocks.DatabaseMock
	service            *CurrencyService
}

func (s *CurrencyServiceTestSuite) SetupTest() {
	s.constants = bootstrap.NewConstants()
	s.currencyRepository = mocks.NewCurrencyRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewCurrencyService(CurrencyServiceDeps{
		Constants:          s.constants,
		CurrencyRepository: s.currencyRepository,
		DB:                 s.db,
	})
}

func (s *CurrencyServiceTestSuite) TestFindCurrencyByID_Success() {
	currency := &entity.Currency{
		Model:       dbmodel.Model{ID: 1},
		Name:        "Iranian Rial",
		Code:        "IRR",
		ConvertRate: 1,
	}
	s.currencyRepository.On("FindCurrencyByID", s.db, uint(1)).Return(currency, nil).Once()

	result, err := s.service.FindCurrencyByID(1)

	s.NoError(err)
	s.NotNil(result)
	s.Equal("IRR", result.Code)
	s.currencyRepository.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestFindCurrencyByID_NotFound() {
	var nilCurrency *entity.Currency
	s.currencyRepository.On("FindCurrencyByID", s.db, uint(99)).Return(nilCurrency, nil).Once()

	result, err := s.service.FindCurrencyByID(99)

	s.Nil(result)
	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.currencyRepository.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestFindCurrencyByID_RepoError() {
	var nilCurrency *entity.Currency
	repoErr := errors.New("db error")
	s.currencyRepository.On("FindCurrencyByID", s.db, uint(1)).Return(nilCurrency, repoErr).Once()

	result, err := s.service.FindCurrencyByID(1)

	s.Nil(result)
	s.ErrorIs(err, repoErr)
	s.currencyRepository.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestFindCurrencyByCode_Success() {
	currency := &entity.Currency{
		Model:       dbmodel.Model{ID: 2},
		Name:        "US Dollar",
		Code:        "USD",
		ConvertRate: 500000,
	}
	s.currencyRepository.On("FindCurrencyByCode", s.db, "USD").Return(currency, nil).Once()

	result, err := s.service.FindCurrencyByCode("USD")

	s.NoError(err)
	s.NotNil(result)
	s.Equal(uint(500000), result.ConvertRate)
	s.currencyRepository.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestFindCurrencyByCode_NotFound() {
	var nilCurrency *entity.Currency
	s.currencyRepository.On("FindCurrencyByCode", s.db, "XYZ").Return(nilCurrency, nil).Once()

	result, err := s.service.FindCurrencyByCode("XYZ")

	s.Nil(result)
	s.Error(err)
	s.currencyRepository.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestGetCurrencies_Success() {
	currencies := []*entity.Currency{
		{Model: dbmodel.Model{ID: 1}, Code: "IRR", ConvertRate: 1},
		{Model: dbmodel.Model{ID: 2}, Code: "USD", ConvertRate: 500000},
	}
	s.currencyRepository.On("GetCurrencies", s.db).Return(currencies, nil).Once()

	results, err := s.service.GetCurrencies()

	s.NoError(err)
	s.Len(results, 2)
	s.currencyRepository.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestGetCurrencies_Empty() {
	s.currencyRepository.On("GetCurrencies", s.db).Return([]*entity.Currency{}, nil).Once()

	results, err := s.service.GetCurrencies()

	s.NoError(err)
	s.Empty(results)
	s.currencyRepository.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestParseCurrency_FieldMapping() {
	currency := entity.Currency{
		Model:       dbmodel.Model{ID: 3},
		Name:        "Euro",
		Code:        "EUR",
		ConvertRate: 600000,
	}

	result, err := s.service.ParseCurrency(currency)

	s.NoError(err)
	s.Equal(currencydto.CurrencyCredential{
		ID:          3,
		Name:        "Euro",
		Code:        "EUR",
		ConvertRate: 600000,
	}, result)
}

func (s *CurrencyServiceTestSuite) TestUpdateCurrency_NotFound() {
	var nilCurrency *entity.Currency
	s.currencyRepository.On("FindCurrencyByID", s.db, uint(99)).Return(nilCurrency, nil).Once()

	name := "Updated"
	err := s.service.UpdateCurrency(currencydto.UpdateCurrencyRequest{ID: 99, Name: &name})

	s.Error(err)
	var notFound exception.NotFoundError
	s.True(errors.As(err, &notFound))
	s.currencyRepository.AssertExpectations(s.T())
}

func (s *CurrencyServiceTestSuite) TestUpdateCurrency_Success() {
	currency := &entity.Currency{
		Model:       dbmodel.Model{ID: 1},
		Name:        "Iranian Rial",
		Code:        "IRR",
		ConvertRate: 1,
	}
	s.currencyRepository.On("FindCurrencyByID", s.db, uint(1)).Return(currency, nil).Once()
	s.db.On("WithTransaction", mock.AnythingOfType("func(database.Database) error")).Return(nil).Once()
	s.currencyRepository.On("UpdateCurrency", s.db, mock.AnythingOfType("entity.Currency")).Return(nil).Once()

	newName := "ریال ایران"
	err := s.service.UpdateCurrency(currencydto.UpdateCurrencyRequest{ID: 1, Name: &newName})

	s.NoError(err)
	s.currencyRepository.AssertExpectations(s.T())
	s.db.AssertExpectations(s.T())
}

func TestCurrencyServiceSuite(t *testing.T) {
	suite.Run(t, new(CurrencyServiceTestSuite))
}
