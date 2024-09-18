package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Config struct {
	Env Environment
}

func (e *Config) LoadEnv(filenames ...string) *Environment {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	e.Env = Environment{
		PostgresConnection: getEnv("POSTGRES_CONNECTION", "postgres connection is required"),
		JWTSecret:          getEnv("JWT_SECRET", "jwt secret is required"),
		MailFrom:           getEnv("MAIL_FROM", "mail from is required"),
		MailPassword:       getEnv("MAIL_PASSWORD", "mail password is required"),
		MailHost:           getEnv("MAIL_HOST", "mail host is required"),
		MailPort:           uint(getEnvInt("MAIL_PORT", 0)),
	}

	return &e.Env
}
