package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type UserRepository interface {
	FindUsers(db database.Database) ([]*entity.User, error)
	FindUserByID(db database.Database, id uint) (*entity.User, error)
	FindUserByPhone(db database.Database, phone string) (*entity.User, error)
	FindUserByEmail(db database.Database, email string) (*entity.User, error)
	CreateUser(db database.Database, user *entity.User) error
	DeleteUserByPhone(db database.Database, phone string) error
	UpdateUser(db database.Database, user *entity.User) error
}
