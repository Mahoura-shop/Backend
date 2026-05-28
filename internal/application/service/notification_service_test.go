package service

import (
	"errors"
	"testing"

	notificationdto "github.com/Mahoura-shop/Backend/internal/application/dto/notification"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/suite"
)

type NotificationServiceTestSuite struct {
	suite.Suite
	repo    *mocks.NotificationRepositoryMock
	db      *mocks.DatabaseMock
	service *NotificationService
}

func (s *NotificationServiceTestSuite) SetupTest() {
	s.repo = mocks.NewNotificationRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewNotificationService(NotificationServiceDeps{
		NotificationRepository: s.repo,
		DB:                     s.db,
	})
}

func (s *NotificationServiceTestSuite) TestCreateNotification_Success() {
	n := entity.Notification{
		UserID: 1,
		Type:   enum.NotificationType(1),
		Title:  "Order Shipped",
		Body:   "Your order is on the way",
	}
	created := &entity.Notification{Type: n.Type, Title: n.Title}
	s.repo.On("Create", s.db, n).Return(created, nil).Once()

	err := s.service.CreateNotification(1, 1, "Order Shipped", "Your order is on the way", nil)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *NotificationServiceTestSuite) TestCreateNotification_RepoError() {
	n := entity.Notification{
		UserID: 1,
		Type:   enum.NotificationType(1),
		Title:  "Test",
		Body:   "Body",
	}
	s.repo.On("Create", s.db, n).Return(nil, errors.New("db error")).Once()

	err := s.service.CreateNotification(1, 1, "Test", "Body", nil)

	s.Error(err)
	s.repo.AssertExpectations(s.T())
}

func (s *NotificationServiceTestSuite) TestGetNotifications_ReturnsCredentials() {
	n := &entity.Notification{
		Type:  enum.NotificationType(1),
		Title: "Test",
		Body:  "Hello",
	}
	s.repo.On("GetByUserID", s.db, uint(1)).Return([]*entity.Notification{n}, nil).Once()

	result, err := s.service.GetNotifications(1)

	s.NoError(err)
	s.Len(result, 1)
	s.Equal("Test", result[0].Title)
	s.repo.AssertExpectations(s.T())
}

func (s *NotificationServiceTestSuite) TestGetNotifications_RepoError() {
	s.repo.On("GetByUserID", s.db, uint(1)).Return(nil, errors.New("db error")).Once()

	result, err := s.service.GetNotifications(1)

	s.Nil(result)
	s.Error(err)
	s.repo.AssertExpectations(s.T())
}

func (s *NotificationServiceTestSuite) TestGetNotifications_Empty() {
	s.repo.On("GetByUserID", s.db, uint(2)).Return([]*entity.Notification{}, nil).Once()

	result, err := s.service.GetNotifications(2)

	s.NoError(err)
	s.Empty(result)
	s.repo.AssertExpectations(s.T())
}

func (s *NotificationServiceTestSuite) TestMarkAsRead_CallsRepo() {
	s.repo.On("MarkAsRead", s.db, uint(5), uint(1)).Return(nil).Once()

	err := s.service.MarkAsRead(5, 1)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *NotificationServiceTestSuite) TestMarkAsRead_RepoError() {
	s.repo.On("MarkAsRead", s.db, uint(5), uint(1)).Return(errors.New("not found")).Once()

	err := s.service.MarkAsRead(5, 1)

	s.Error(err)
	s.repo.AssertExpectations(s.T())
}

func (s *NotificationServiceTestSuite) TestMarkAllAsRead_CallsRepo() {
	s.repo.On("MarkAllAsRead", s.db, uint(1)).Return(nil).Once()

	err := s.service.MarkAllAsRead(1)

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *NotificationServiceTestSuite) TestCountUnread_ReturnsCount() {
	s.repo.On("CountUnread", s.db, uint(1)).Return(3, nil).Once()

	count, err := s.service.CountUnread(1)

	s.NoError(err)
	s.Equal(int64(3), count)
	s.repo.AssertExpectations(s.T())
}

func (s *NotificationServiceTestSuite) TestGetNotifications_MapsFields() {
	refID := uint(99)
	n := &entity.Notification{
		Type:   enum.NotificationType(2),
		Title:  "Order Update",
		Body:   "Your order shipped",
		IsRead: true,
		RefID:  &refID,
	}
	s.repo.On("GetByUserID", s.db, uint(7)).Return([]*entity.Notification{n}, nil).Once()

	result, err := s.service.GetNotifications(7)

	s.NoError(err)
	s.Require().Len(result, 1)
	cred := result[0]
	s.Equal("Order Update", cred.Title)
	s.Equal("Your order shipped", cred.Body)
	s.True(cred.IsRead)
	s.Require().NotNil(cred.RefID)
	s.Equal(uint(99), *cred.RefID)
	s.repo.AssertExpectations(s.T())
}

func (s *NotificationServiceTestSuite) TestGetNotifications_ReturnsCorrectType() {
	n := &entity.Notification{
		Type:  enum.NotificationType(1),
		Title: "Test",
		Body:  "Body",
	}
	s.repo.On("GetByUserID", s.db, uint(1)).Return([]*entity.Notification{n}, nil).Once()

	result, err := s.service.GetNotifications(1)

	s.NoError(err)
	_, ok := interface{}(result).([]notificationdto.NotificationCredential)
	s.True(ok)
	s.repo.AssertExpectations(s.T())
}

func TestNotificationServiceSuite(t *testing.T) {
	suite.Run(t, new(NotificationServiceTestSuite))
}
