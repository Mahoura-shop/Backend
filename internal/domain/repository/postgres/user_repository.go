package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type UserRepository interface {
	FindUsers(database.Database) ([]*entity.User, error)
	FindUserByID(database.Database, uint) (*entity.User, error)
	FindUserByPhone(database.Database, string) (*entity.User, error)
	FindUserByEmail(database.Database, string) (*entity.User, error)
	CreateUser(database.Database, *entity.User) error
	DeleteUserByPhone(database.Database, string) error
	UpdateUser(database.Database, *entity.User) error
}
