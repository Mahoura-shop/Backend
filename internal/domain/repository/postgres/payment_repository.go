package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type PaymentRepository interface {
	CreatePayment(db database.Database, payment entity.Payment) (*entity.Payment, error)
	FindPaymentByAuthority(db database.Database, authority string) (*entity.Payment, error)
	FindPaymentByOrderID(db database.Database, orderID uint) (*entity.Payment, error)
	UpdatePayment(db database.Database, payment entity.Payment) error
}
