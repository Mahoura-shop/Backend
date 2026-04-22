package seed

import (
	"fmt"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	repository "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

var currencies = []entity.Currency{
	{
		Name:        "ریال",
		Code:        "IRR",
		ConvertRate: 1,
	},
	{
		Name:        "درهم",
		Code:        "AED",
		ConvertRate: 450000,
	},
	{
		Name:        "دلار",
		Code:        "USD",
		ConvertRate: 1650000,
	},
}


type CurrencySeeder struct {
	currencyRepository repository.CurrencyRepository
	db                database.Database
}

func NewCurrencySeeder(
	currencyRepository repository.CurrencyRepository,
	db database.Database,
) *CurrencySeeder {
	return &CurrencySeeder{
		currencyRepository: currencyRepository,
		db:                 db,
	}
}

func (seeder *CurrencySeeder) SeedCurrencies() {
	for _, currency := range currencies {
		existing, err := seeder.currencyRepository.FindCurrencyByCode(seeder.db, currency.Code)
		if err != nil {
			panic(err)
		}

		if existing != nil {
			continue
		}

		err = seeder.currencyRepository.CreateCurrency(seeder.db, currency)
		if err != nil {
			panic(fmt.Errorf("error creating currency %s: %w", currency.Code, err))
		}
	}
}
