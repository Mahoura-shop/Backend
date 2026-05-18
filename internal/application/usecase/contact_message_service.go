package usecase

import "github.com/Mahoura-shop/Backend/internal/domain/entity"

type ContactMessageService interface {
	Submit(name, email, subject, message string) error
	GetAll() ([]*entity.ContactMessage, error)
}
