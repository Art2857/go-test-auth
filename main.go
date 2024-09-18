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

	database.InitDB(config.Env.POSTGRES_CONNECTION)
	defer database.CloseDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/docs/*", echoSwagger.WrapHandler)

	mailService := mail.MailService{
		From:     config.Env.MAIL_FROM,
		Password: config.Env.MAIL_PASSWORD,
		Host:     config.Env.MAIL_HOST,
		Port:     config.Env.MAIL_PORT,
	}

	tokenRepository := token.NewRepository(database.DB)
	tokenService := token.NewService(tokenRepository, &mailService)

	tokenHandlers := token.NewHanders(tokenService)
	tokenHandlers.InitHandlers(e)

	log.Println("Server is running on port 8000...")
	log.Fatal(e.Start(":8000"))
}
