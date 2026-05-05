package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type UpgradeRequestRepository interface {
	CreateUpgradeRequest(database.Database, entity.UpgradeRequest) (*entity.UpgradeRequest, error)
	FindUpgradeRequestByID(database.Database, uint) (*entity.UpgradeRequest, error)
	GetUpgradeRequests(database.Database) ([]*entity.UpgradeRequest, error)
	GetUpgradeRequestsByUserID(database.Database, uint) ([]*entity.UpgradeRequest, error)
	GetUpgradeRequestsByStatus(database.Database, enum.UpgradeRequestStatus) ([]*entity.UpgradeRequest, error)
	FindPendingByUserID(database.Database, uint) (*entity.UpgradeRequest, error)
	UpdateUpgradeRequest(database.Database, entity.UpgradeRequest) error
}
