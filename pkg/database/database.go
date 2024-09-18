package database

import (
	"auth-service/pkg/database/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB инициализирует подключение к базе данных и миграцию
func InitDB(connectString string) {
	var err error

	DB, err = gorm.Open(postgres.Open(connectString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	log.Println("Connected to the database.")

	InitMigration()
}

func isDatabaseEmpty(db *gorm.DB) (bool, error) {
	var count int64
	err := db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';").Scan(&count).Error
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func InitMigration() {
	isEmpty, err := isDatabaseEmpty(DB)
	if err != nil {
		log.Fatal(err)
	}

	if isEmpty {
		log.Println("Database is empty, running init migration...")

		err := DB.AutoMigrate(&models.RefreshTokenModel{})
		if err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}

		log.Println("Init Migration complete.")
	}
}

func CloseDB() {
	db, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database object: %v", err)
	}

	err = db.Close()
	if err != nil {
		log.Fatalf("Error closing the database connection: %v", err)
	}

	log.Println("Database connection closed.")
}
