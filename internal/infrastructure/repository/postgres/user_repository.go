package postgres

import (
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	repository "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
	"gorm.io/gorm"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (repo *UserRepository) FindUsers(db database.Database) ([]*entity.User, error) {
	var users []*entity.User
	result := db.GetDB().Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (repo *UserRepository) FindUserByID(db database.Database, id uint) (*entity.User, error) {
	var user entity.User
	result := db.GetDB().First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) FindUserByStatus(db database.Database, statuses []enum.UserStatus, opts ...repository.QueryModifier) ([]*entity.User, error) {
	var users []*entity.User
	query := db.GetDB().Where("status IN ?", statuses)

	for _, opt := range opts {
		query = opt.Apply(query).(*gorm.DB)
	}

	result := query.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (repo *UserRepository) FindUserByEmail(db database.Database, email string) (*entity.User, error) {
	var user entity.User
	result := db.GetDB().Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) FindUserByPhone(db database.Database, phone string) (*entity.User, error) {
	var user entity.User
	result := db.GetDB().Where("phone = ?", phone).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) CreateUser(db database.Database, user entity.User) (error) {
	return db.GetDB().Create(&user).Error
}

func (repo *UserRepository) DeleteUserByPhone(db database.Database, phone string) (error) {
	return db.GetDB().Where("phone = ?", phone).Unscoped().Delete(&entity.User{}).Error
}

func (repo *UserRepository) UpdateUser(db database.Database, user entity.User) (error) {
	return db.GetDB().Save(&user).Error
}