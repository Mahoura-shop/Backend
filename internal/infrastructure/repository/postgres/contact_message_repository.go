package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type ContactMessageRepository struct{}

func NewContactMessageRepository() *ContactMessageRepository {
	return &ContactMessageRepository{}
}

func (r *ContactMessageRepository) Create(db database.Database, msg entity.ContactMessage) (*entity.ContactMessage, error) {
	result := db.GetDB().Create(&msg)
	if result.Error != nil {
		return nil, result.Error
	}
	return &msg, nil
}

func (r *ContactMessageRepository) GetAll(db database.Database) ([]*entity.ContactMessage, error) {
	msgs := make([]*entity.ContactMessage, 0)
	result := db.GetDB().Order("created_at desc").Find(&msgs)
	if result.Error != nil {
		return nil, result.Error
	}
	return msgs, nil
}
