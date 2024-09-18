package main

import (
	"auth-service/pkg/config"
	"auth-service/pkg/database"
	"auth-service/pkg/mail"
	"auth-service/pkg/token"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "auth-service/docs"
)

func main() {
	config.Init()

	database.InitDB(config.Env.PostgresConnection)
	defer database.CloseDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/docs/*", echoSwagger.WrapHandler)

	mailService := mail.MailService{
		From:     config.Env.MailFrom,
		Password: config.Env.MailPassword,
		Host:     config.Env.MailHost,
		Port:     config.Env.MailPort,
	}

	tokenRepository := token.NewRepository(database.DB)
	tokenService := token.NewService(tokenRepository, &mailService)

	tokenHandlers := token.NewHanders(tokenService)
	tokenHandlers.InitHandlers(e)

	log.Println("Server is running on port 8000...")
	log.Fatal(e.Start(":8000"))
}
