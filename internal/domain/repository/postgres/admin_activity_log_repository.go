package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type AdminActivityLogRepository interface {
	Create(db database.Database, log entity.AdminActivityLog) error
	GetAll(db database.Database) ([]*entity.AdminActivityLog, error)
}
