package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type CurrencyRepository interface {
	FindCurrencyByID(database.Database, uint) (*entity.Currency, error)
	FindCurrencyByCode(database.Database, string) (*entity.Currency, error)
	GetCurrencies(database.Database) ([]*entity.Currency, error)
	UpdateCurrency(database.Database, *entity.Currency) error
}
