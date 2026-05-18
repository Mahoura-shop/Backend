package main

import (
	"fmt"
	"log"

	"github.com/Mahoura-shop/Backend/bootstrap"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	godotenv.Load(".env")
	config := bootstrap.Run()

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		config.Env.PrimaryDB.Host,
		config.Env.PrimaryDB.User,
		config.Env.PrimaryDB.Password,
		config.Env.PrimaryDB.Name,
		config.Env.PrimaryDB.Port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	result := db.Exec("UPDATE users SET type = 4 WHERE type = 3")
	if result.Error != nil {
		log.Fatalf("migration failed: %v", result.Error)
	}

	log.Printf("migration complete: %d user(s) migrated from shopkeeperCheque to shopkeeper", result.RowsAffected)
}
