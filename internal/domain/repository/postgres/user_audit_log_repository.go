package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type UserAuditLogRepository interface {
	CreateUserAuditLog(database.Database, entity.UserAuditLog) error
	GetAuditLogsByUserID(database.Database, uint) ([]*entity.UserAuditLog, error)
}
