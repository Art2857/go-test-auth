package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	Postgres Postgres
}

func (db *DB) OpenPostgres(connectString string) *Postgres {
	var err error

	postgres, err := gorm.Open(postgres.Open(connectString), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	log.Println("Connected to the database.")

	db.Postgres = Postgres{postgres}

	return &db.Postgres
}

func (db *DB) ClosePostgres() {
	sqlDB, err := db.Postgres.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get database object: %v", err)
	}

	err = sqlDB.Close()
	if err != nil {
		log.Fatalf("Error closing the database connection: %v", err)
	}

	log.Println("Database connection closed.")
}
