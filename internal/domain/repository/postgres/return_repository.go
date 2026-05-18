package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type ReturnRepository interface {
	CreateReturn(database.Database, entity.Return) (*entity.Return, error)
	FindReturnByID(database.Database, uint) (*entity.Return, error)
	GetReturnsByUserID(database.Database, uint) ([]*entity.Return, error)
	GetAllReturns(database.Database) ([]*entity.Return, error)
	GetReturnsByStatus(database.Database, enum.ReturnStatus) ([]*entity.Return, error)
	UpdateReturn(database.Database, entity.Return) error
}
