package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type VisitsByDay struct {
	Date  string
	Count uint
}

type ProductVisitRepository interface {
	CreateVisit(database.Database, entity.ProductVisit) (*entity.ProductVisit, error)
	GetVisitCountByProductID(database.Database, uint) (int64, error)
	HasVisitedInLast24h(database.Database, uint, string) (bool, error)
	GetVisitsPerDay(database.Database, uint, int) ([]VisitsByDay, error)
	GetAllVisitsPerDay(database.Database, int) ([]VisitsByDay, error)
}
