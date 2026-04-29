package postgres

import (
	"errors"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type PaymentRepository struct{}

func NewPaymentRepository() *PaymentRepository {
	return &PaymentRepository{}
}

func (repo *PaymentRepository) CreatePayment(db database.Database, payment entity.Payment) (*entity.Payment, error) {
	result := db.GetDB().Create(&payment)
	if result.Error != nil {
		return nil, result.Error
	}
	return &payment, nil
}

func (repo *PaymentRepository) FindPaymentByAuthority(db database.Database, authority string) (*entity.Payment, error) {
	var payment entity.Payment
	result := db.GetDB().Where("authority = ?", authority).First(&payment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &payment, nil
}

func (repo *PaymentRepository) FindPaymentByOrderID(db database.Database, orderID uint) (*entity.Payment, error) {
	var payment entity.Payment
	result := db.GetDB().Where("order_id = ?", orderID).First(&payment)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	return &payment, nil
}

func (repo *PaymentRepository) UpdatePayment(db database.Database, payment entity.Payment) error {
	return db.GetDB().Save(&payment).Error
}
