package test_utils

import (
	"auth/pkg/config"
	"auth/pkg/database"
	"auth/pkg/mail"
	"auth/pkg/token"
	"log"

	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func ClearDB(db *gorm.DB) {
	db.Transaction(func(tx *gorm.DB) error {
		tables, err := tx.Migrator().GetTables()
		if err != nil {
			return err
		}

		for _, table := range tables {
			err := tx.Migrator().DropTable(table)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

type Application struct {
	Router        *echo.Echo
	ConfigService *config.Config
	DB            *database.DB

	MailService *mail.MailService

	TokenRepository *token.TokenRepository
	TokenService    *token.TokenService
	TokenHandlers   *token.TokenHandlers
}

func SetupTest() *Application {
	// Config
	config := config.Config{}
	config.LoadEnv("../.env.test")

	// Database
	db := database.DB{}

	// Postgres
	db.OpenPostgres(config.Env.PostgresConnection)
	// defer db.ClosePostgres()

	if empty, _ := db.Postgres.IsDBEmpty(); !empty {
		log.Fatal("Database is not empty")
	}

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

	tokenHandlers := token.NewHanders(tokenService)
	tokenHandlers.InitHandlers(e)

	return &Application{
		Router:          e,
		ConfigService:   &config,
		DB:              &db,
		MailService:     &mailService,
		TokenRepository: tokenRepository,
		TokenService:    tokenService,
		TokenHandlers:   tokenHandlers,
	}
}

func EndTest(app *Application) {
	ClearDB(app.DB.Postgres.DB)

	app.DB.ClosePostgres()

	if err := app.Router.Shutdown(context.Background()); err != nil {
		app.Router.Logger.Fatal(err)
	}
}
