package main

import (
	"auth/pkg/config"
	"auth/pkg/database"
	"auth/pkg/mail"
	"auth/pkg/token"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "auth/docs"
)

func main() {
	// Config
	config := config.Config{}
	config.LoadEnv()

	// Database
	db := database.DB{}

	// Postgres
	db.OpenPostgres(config.Env.PostgresConnection)
	defer db.ClosePostgres()

	db.Postgres.InitMigration()

	// Services
	mailService := mail.MailService{
		From:     config.Env.MailFrom,
		Password: config.Env.MailPassword,
		Host:     config.Env.MailHost,
		Port:     config.Env.MailPort,
	}

	tokenRepository := token.NewRepository(&db.Postgres)
	tokenService := token.NewService(&config, tokenRepository, &mailService)

	// Handlers
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/docs/*", echoSwagger.WrapHandler)

	tokenHandlers := token.NewHanders(tokenService)
	tokenHandlers.InitHandlers(e)

	log.Println("Server is running on port 8000...")
	log.Fatal(e.Start(":8000"))
}
