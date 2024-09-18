package main

import (
	"auth-service/pkg/config"
	"auth-service/pkg/database"
	"auth-service/pkg/mail"
	"auth-service/pkg/token"
	"log"
	"os"

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
		From:     os.Getenv("MAIL_FROM"),
		Password: os.Getenv("MAIL_PASSWORD"),
		Host:     os.Getenv("MAIL_HOST"),
		Port:     os.Getenv("MAIL_PORT"),
	}

	tokenRepository := token.NewRepository(database.DB)
	tokenService := token.NewService(tokenRepository, &mailService)

	tokenHandlers := token.NewHanders(tokenService)
	tokenHandlers.InitHandlers(e)

	log.Println("Server is running on port 8000...")
	log.Fatal(e.Start(":8000"))
}
