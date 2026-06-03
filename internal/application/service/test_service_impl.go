package service

import (
	"os"

	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/Mahoura-shop/Backend/internal/infrastructure/database"
)

type TestService struct {
	constants *bootstrap.Constants
	admins    *bootstrap.AdminCredentials
	db        database.Database
}

func NewTestService(
	constants *bootstrap.Constants,
	admins *bootstrap.AdminCredentials,
	db database.Database,
) *TestService {
	return &TestService{
		constants: constants,
		admins:    admins,
		db:        db,
	}
}

func (s *TestService) Test(test string) (string, error) {
	return "This is " + test, nil
}

func (s *TestService) ResetDB() error {
	if os.Getenv("APP_MODE") != "test" {
		return nil
	}

	truncateSQL := `
		DO $$
		DECLARE r RECORD;
		BEGIN
			FOR r IN (
				SELECT tablename FROM pg_tables
				WHERE schemaname = 'public'
				AND tablename NOT IN (
					'provinces', 'cities', 'currencies',
					'permissions', 'roles', 'role_permissions'
				)
			) LOOP
				EXECUTE 'TRUNCATE TABLE "' || r.tablename || '" RESTART IDENTITY CASCADE';
			END LOOP;
		END $$;
	`
	if err := s.db.GetDB().Exec(truncateSQL).Error; err != nil {
		return err
	}

	for _, admin := range s.admins.Admins {
		insertSQL := `
			INSERT INTO users (created_at, updated_at, phone, status, is_admin)
			VALUES (NOW(), NOW(), ?, 1, true)
			ON CONFLICT (phone) DO UPDATE SET is_admin = true, status = 1
		`
		if err := s.db.GetDB().Exec(insertSQL, admin.Phone).Error; err != nil {
			return err
		}
	}

	return nil
}
