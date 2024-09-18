package database

import (
	"auth/pkg/database/models"
	"log"

	"gorm.io/gorm"
)

type Postgres struct {
	*gorm.DB
}

func (postgres *Postgres) IsDBEmpty() (bool, error) {
	tables, err := postgres.Migrator().GetTables()
	if err != nil {
		return false, err
	}

	return len(tables) == 0, nil
}

func (postgres Postgres) InitMigration() {
	isEmpty, err := postgres.IsDBEmpty()
	if err != nil {
		log.Fatal(err)
	}

	if isEmpty {
		log.Println("Database is empty, running init migration...")

		err := postgres.AutoMigrate(&models.RefreshTokenModel{})
		if err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}

		log.Println("Init Migration complete.")
	}
}
