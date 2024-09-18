package config

import (
	"log"
	"os"
	"strconv"
)

type Environment struct {
	PostgresConnection string
	JWTSecret          string
	MailFrom           string
	MailPassword       string
	MailHost           string
	MailPort           uint
}

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
