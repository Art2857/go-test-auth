package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v4"
)

var DB *pgx.Conn

func InitDB(connStr string) {
	var err error

	DB, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	log.Println("Connected to the database.")

	InitMigration()
}

func isDatabaseEmpty(db *pgx.Conn) (bool, error) {
	var count int
	err := db.QueryRow(context.Background(), "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public';").Scan(&count)
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

		// По хорошему добавить что-то вроде expires_at, но в случае redis'а это было бы проще
		_, err := DB.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS refresh_tokens (
			id SERIAL PRIMARY KEY,
			token_hash TEXT NOT NULL,
			ip TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT NOW()
		);`)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Init Migration complete.")
	}
}

func CloseDB() {
	err := DB.Close(context.Background())
	if err != nil {
		log.Fatalf("Error closing the database connection: %v", err)
	}

	log.Println("Database connection closed.")
}
