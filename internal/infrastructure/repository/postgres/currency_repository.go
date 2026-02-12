package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type CurrencyRepository struct{}

func NewCurrencyRepository() *CurrencyRepository {
	return &CurrencyRepository{}
}

func (repo *CurrencyRepository) FindCurrencyByID(db database.Database, currencyID uint) (*entity.Currency, error) {
	var currency entity.Currency
	result := db.GetDB().Where("id = ?", currencyID).First(&currency)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &currency, nil
}

func (repo *CurrencyRepository) FindCurrencyByCode(db database.Database, code string) (*entity.Currency, error) {
	var currency entity.Currency
	result := db.GetDB().Where("code = ?", code).First(&currency)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &currency, nil
}

func (repo *CurrencyRepository) CreateCurrency(db database.Database, currency *entity.Currency) error {
	result := db.GetDB().Create(&currency)
	if result.Error != nil {
        return result.Error
    }
    
    return nil
}

func (repo *CurrencyRepository) GetCurrencies(db database.Database) ([]*entity.Currency, error) {
	var currencies []*entity.Currency
	result := db.GetDB().Find(&currencies)
	
	if result.Error != nil {
		return nil, result.Error
	}
	
	return currencies, nil
}

func (repo *CurrencyRepository) UpdateCurrency(db database.Database, currency *entity.Currency) error {
	return db.GetDB().Save(&currency).Error
}