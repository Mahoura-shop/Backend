package mocks

import (
	returndto "github.com/Mahoura-shop/Backend/internal/application/dto/return"
	"github.com/stretchr/testify/mock"
)

type ReturnServiceMock struct {
	mock.Mock
}

func NewReturnServiceMock() *ReturnServiceMock {
	return &ReturnServiceMock{}
}

func (m *ReturnServiceMock) RequestReturn(req returndto.RequestReturnRequest) error {
	args := m.Called(req)
	return args.Error(0)
}

func (m *ReturnServiceMock) GetMyReturns(userID uint) ([]returndto.ReturnCredential, error) {
	args := m.Called(userID)
	return args.Get(0).([]returndto.ReturnCredential), args.Error(1)
}

func (m *ReturnServiceMock) GetAllReturns(status string) ([]returndto.ReturnCredential, error) {
	args := m.Called(status)
	return args.Get(0).([]returndto.ReturnCredential), args.Error(1)
}

func (m *ReturnServiceMock) ReviewReturn(returnID uint, req returndto.ReviewReturnRequest) error {
	args := m.Called(returnID, req)
	return args.Error(0)
}

func (m *ReturnServiceMock) ProcessRefund(returnID uint, adminID uint) error {
	args := m.Called(returnID, adminID)
	return args.Error(0)
}
