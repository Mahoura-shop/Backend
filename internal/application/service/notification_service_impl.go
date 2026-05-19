package service

import (
	notificationdto "github.com/Mahoura-shop/Backend/internal/application/dto/notification"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	domainPostgres "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/sse"
)

type NotificationService struct {
	notificationRepository domainPostgres.NotificationRepository
	db                     database.Database
}

type NotificationServiceDeps struct {
	NotificationRepository domainPostgres.NotificationRepository
	DB                     database.Database
}

func NewNotificationService(deps NotificationServiceDeps) *NotificationService {
	return &NotificationService{
		notificationRepository: deps.NotificationRepository,
		db:                     deps.DB,
	}
}

func (s *NotificationService) CreateNotification(userID uint, notifType uint, title, body string, refID *uint) error {
	n := entity.Notification{
		UserID: userID,
		Type:   enum.NotificationType(notifType),
		Title:  title,
		Body:   body,
		RefID:  refID,
	}
	created, err := s.notificationRepository.Create(s.db, n)
	if err != nil {
		return err
	}
	sse.Global.Send(userID, sse.Event{
		ID:    created.ID,
		Type:  notifType,
		Title: title,
		Body:  body,
		RefID: refID,
	})
	return nil
}

func (s *NotificationService) GetNotifications(userID uint) ([]notificationdto.NotificationCredential, error) {
	notifications, err := s.notificationRepository.GetByUserID(s.db, userID)
	if err != nil {
		return nil, err
	}
	var result []notificationdto.NotificationCredential
	for _, n := range notifications {
		result = append(result, notificationdto.NotificationCredential{
			ID:        n.ID,
			Type:      n.Type,
			Title:     n.Title,
			Body:      n.Body,
			IsRead:    n.IsRead,
			RefID:     n.RefID,
			CreatedAt: n.CreatedAt,
		})
	}
	return result, nil
}

func (s *NotificationService) MarkAsRead(id, userID uint) error {
	return s.notificationRepository.MarkAsRead(s.db, id, userID)
}

func (s *NotificationService) MarkAllAsRead(userID uint) error {
	return s.notificationRepository.MarkAllAsRead(s.db, userID)
}

func (s *NotificationService) CountUnread(userID uint) (int64, error) {
	return s.notificationRepository.CountUnread(s.db, userID)
}
