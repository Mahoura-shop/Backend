package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type ContactMessageRepository interface {
	Create(db database.Database, msg entity.ContactMessage) (*entity.ContactMessage, error)
	GetAll(db database.Database) ([]*entity.ContactMessage, error)
}
