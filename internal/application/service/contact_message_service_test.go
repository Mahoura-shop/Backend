package service

import (
	"errors"
	"testing"

	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/mocks"
	"github.com/stretchr/testify/suite"
)

type ContactMessageServiceTestSuite struct {
	suite.Suite
	repo    *mocks.ContactMessageRepositoryMock
	db      *mocks.DatabaseMock
	service *ContactMessageService
}

func (s *ContactMessageServiceTestSuite) SetupTest() {
	s.repo = mocks.NewContactMessageRepositoryMock()
	s.db = mocks.NewDatabaseMock()
	s.service = NewContactMessageService(ContactMessageServiceDeps{
		ContactMessageRepository: s.repo,
		DB:                       s.db,
	})
}

func (s *ContactMessageServiceTestSuite) TestSubmit_Success() {
	msg := entity.ContactMessage{
		Name:    "Ali",
		Email:   "ali@example.com",
		Subject: "Question",
		Message: "Hello there",
	}
	created := &entity.ContactMessage{Name: "Ali"}
	s.repo.On("Create", s.db, msg).Return(created, nil).Once()

	err := s.service.Submit("Ali", "ali@example.com", "Question", "Hello there")

	s.NoError(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ContactMessageServiceTestSuite) TestSubmit_RepoError() {
	msg := entity.ContactMessage{
		Name:    "Sara",
		Email:   "sara@example.com",
		Subject: "Bug",
		Message: "Found an issue",
	}
	s.repo.On("Create", s.db, msg).Return(nil, errors.New("db error")).Once()

	err := s.service.Submit("Sara", "sara@example.com", "Bug", "Found an issue")

	s.Error(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ContactMessageServiceTestSuite) TestGetAll_ReturnsList() {
	msgs := []*entity.ContactMessage{
		{Name: "Ali", Email: "ali@example.com"},
		{Name: "Sara", Email: "sara@example.com"},
	}
	s.repo.On("GetAll", s.db).Return(msgs, nil).Once()

	result, err := s.service.GetAll()

	s.NoError(err)
	s.Len(result, 2)
	s.Equal("Ali", result[0].Name)
	s.repo.AssertExpectations(s.T())
}

func (s *ContactMessageServiceTestSuite) TestGetAll_RepoError() {
	s.repo.On("GetAll", s.db).Return(nil, errors.New("db error")).Once()

	result, err := s.service.GetAll()

	s.Nil(result)
	s.Error(err)
	s.repo.AssertExpectations(s.T())
}

func (s *ContactMessageServiceTestSuite) TestGetAll_Empty() {
	s.repo.On("GetAll", s.db).Return([]*entity.ContactMessage{}, nil).Once()

	result, err := s.service.GetAll()

	s.NoError(err)
	s.Empty(result)
	s.repo.AssertExpectations(s.T())
}

func TestContactMessageServiceSuite(t *testing.T) {
	suite.Run(t, new(ContactMessageServiceTestSuite))
}
