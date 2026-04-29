package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type InstalmentRepository interface {
	CreateInstalment(db database.Database, instalment entity.Instalment) error
	GetInstalmentsByOrderID(db database.Database, orderID uint) ([]*entity.Instalment, error)
	UpdateInstalment(db database.Database, instalment entity.Instalment) error
}
