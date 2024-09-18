package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environment struct {
	POSTGRES_CONNECTION string
	JWT_SECRET          string
	MAIL_FROM           string
	MAIL_PASSWORD       string
	MAIL_HOST           string
	MAIL_PORT           string
}

var Env Environment

func getEnv(key string, message string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatal(message)
	}

	return value
}

func Init(filenames ...string) Environment {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	Env = Environment{
		POSTGRES_CONNECTION: getEnv("POSTGRES_CONNECTION", "postgres connection is required"),
		JWT_SECRET:          getEnv("JWT_SECRET", "jwt secret is required"),
		MAIL_FROM:           getEnv("MAIL_FROM", "mail from is required"),
		MAIL_PASSWORD:       getEnv("MAIL_PASSWORD", "mail password is required"),
		MAIL_HOST:           getEnv("MAIL_HOST", "mail host is required"),
		MAIL_PORT:           getEnv("MAIL_PORT", "mail port is required"),
	}

	return Env
}
