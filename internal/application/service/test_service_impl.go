package service

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type TestService struct {
	constants         *bootstrap.Constants
	db                database.Database
}

func NewTestService(
	constants *bootstrap.Constants,
	db database.Database,
) *TestService {
	return &TestService{
		constants:           constants,
		db:                  db,
	}
}

func (testService *TestService) Test(test string) (string, error) {
	result := "This is " + test
	return result, nil
}