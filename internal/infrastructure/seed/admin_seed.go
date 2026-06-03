package seed

import (
	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/domain/entity"
	"github.com/Mahoura-shop/Backend/internal/domain/enum"
	repository "github.com/Mahoura-shop/Backend/internal/domain/repository/postgres"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type AdminSeeder struct {
	admins         *bootstrap.AdminCredentials
	userRepository repository.UserRepository
	db             database.Database
}

func NewAdminSeeder(
	admins *bootstrap.AdminCredentials,
	userRepository repository.UserRepository,
	db database.Database,
) *AdminSeeder {
	return &AdminSeeder{
		admins:         admins,
		userRepository: userRepository,
		db:             db,
	}
}

func (adminSeeder *AdminSeeder) SeedAdmins() {
	for _, admin := range adminSeeder.admins.Admins {
		adminSeeder.getOrCreateAdmin(admin)
	}
}

func (adminSeeder *AdminSeeder) getOrCreateAdmin(admin bootstrap.AdminAccount) *entity.User {
	user, err := adminSeeder.userRepository.FindUserByPhone(adminSeeder.db, admin.Phone)
	if err != nil {
		panic(err)
	}

	if user == nil {
		user = &entity.User{
			Phone:     admin.Phone,
			Status:    enum.UserStatusActive,
			FirstName: admin.FirstName,
			LastName:  admin.LastName,
			IsAdmin:   true,
		}

		if err := adminSeeder.userRepository.CreateUser(adminSeeder.db, *user); err != nil {
			panic(err)
		}

		return user
	}

	user.IsAdmin = true

	if err := adminSeeder.userRepository.UpdateUser(adminSeeder.db, *user); err != nil {
		panic(err)
	}

	return user
}
