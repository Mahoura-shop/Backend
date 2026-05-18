package service

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	domainPostgres "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type ContactMessageService struct {
	repo domainPostgres.ContactMessageRepository
	db   database.Database
}

type ContactMessageServiceDeps struct {
	ContactMessageRepository domainPostgres.ContactMessageRepository
	DB                       database.Database
}

func NewContactMessageService(deps ContactMessageServiceDeps) *ContactMessageService {
	return &ContactMessageService{
		repo: deps.ContactMessageRepository,
		db:   deps.DB,
	}
}

func (s *ContactMessageService) Submit(name, email, subject, message string) error {
	msg := entity.ContactMessage{
		Name:    name,
		Email:   email,
		Subject: subject,
		Message: message,
	}
	_, err := s.repo.Create(s.db, msg)
	return err
}

func (s *ContactMessageService) GetAll() ([]*entity.ContactMessage, error) {
	return s.repo.GetAll(s.db)
}
