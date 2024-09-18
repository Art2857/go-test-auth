package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environment struct {
	PostgresConnection string
	JWTSecret          string
	MailFrom           string
	MailPassword       string
	MailHost           string
	MailPort           uint
}

var Env Environment

func getEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("Environment variable %s is required: %s", key, fallback)
	}
	return value
}

func getEnvInt(key string, defValue int) int {
	valueStr := getEnv(key, "")

	if valueStr == "" {
		return defValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Fatalf("Invalid integer value for environment variable %s: %s", key, valueStr)
	}
	return value
}

func Init(filenames ...string) Environment {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	Env = Environment{
		PostgresConnection: getEnv("POSTGRES_CONNECTION", "postgres connection is required"),
		JWTSecret:          getEnv("JWT_SECRET", "jwt secret is required"),
		MailFrom:           getEnv("MAIL_FROM", "mail from is required"),
		MailPassword:       getEnv("MAIL_PASSWORD", "mail password is required"),
		MailHost:           getEnv("MAIL_HOST", "mail host is required"),
		MailPort:           uint(getEnvInt("MAIL_PORT", 0)),
	}

	return Env
}
