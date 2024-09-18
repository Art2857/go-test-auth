package token_test

import (
	"auth-service/pkg/config"
	"auth-service/pkg/database"
	"auth-service/pkg/mail"
	"auth-service/pkg/token"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Написал только часть теста, так как много чего ещё дописывать нужно, чтобы всё затестить...

func setupRouter() (*echo.Echo, *token.TokenService) {
	config.Init("../.env.test")

	database.InitDB(config.Env.PostgresConnection)
	// defer database.CloseDB()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

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

	return e, tokenService
}

// Кейсы для тестирования:
// Должен вернуть код 200
// Должен вернуть access_token и refresh_token
// Должен вернуть access_token типа jwt и refresh_token типа base64
// Должен вернуть claims(jwt) валидируя пару access и refresh токенов как связанные
// Должен валидировать, что claims(jwt) содержит: userID и является uuid, RefreshTokenHash - является строкой

// Работа с бд...
func TestHandlerGenerateTokenPair(t *testing.T) {
	e, tokenService := setupRouter()

	// Test data
	requestBody := map[string]string{
		"user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
	}
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/token/login", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	// Execute request
	e.ServeHTTP(rec, req)

	// Assert results
	assert.Equal(t, http.StatusOK, rec.Code)

	var tokenPair *token.TokenPair
	err := json.Unmarshal(rec.Body.Bytes(), &tokenPair)
	assert.NoError(t, err)

	accessToken, refreshToken := tokenPair.AccessToken, tokenPair.RefreshToken

	assert.NotEmpty(t, accessToken, "Access token should not be empty")
	assert.NotEmpty(t, refreshToken, "Refresh token should not be empty")

	_, _, err = new(jwt.Parser).ParseUnverified(tokenPair.AccessToken, jwt.MapClaims{})
	assert.NoError(t, err, "Access token should be a valid JWT")

	_, err = base64.StdEncoding.DecodeString(tokenPair.RefreshToken)
	assert.NoError(t, err, "Refresh token should be a valid Base64 encoded string")

	claims, err := tokenService.VerifyTokenPair(tokenPair)
	assert.NoError(t, err, "TokenPair verification should not return an error")
	assert.NotNil(t, claims, "Claims should not be nil")

	assert.NotEmpty(t, claims.UserID, "UserID should not be empty")
	assert.NotEmpty(t, claims.RefreshTokenHash, "RefreshTokenHash should not be empty")

	_, err = uuid.Parse(claims.UserID)
	assert.NoError(t, err, "UserID should be a valid UUID")
	assert.IsType(t, "", claims.RefreshTokenHash, "RefreshTokenHash should be of type string")
}

// При верных данных идёт проверка из первого теста
// При другом ip должен быть отправлено предупреждение на почту
// При неверном access ошибка, при неверном refresh ошибка, если в бд нету записи refresh - ошибка
// func TestHandlerRefreshTokenPair(t *testing.T) {
// 	e, tokenService := setupRouter()

// 	// Test data
// 	requestBody := map[string]string{
// 		"access_token":  "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiYTBlZWJjOTktOWMwYi00ZWY4LWJiNmQtNmJiOWJkMzgwYTExIiwiaXAiOiI6OjEiLCJyZWZyZXNoX3Rva2VuX2hhc2giOiIkMmEkMTAkbDJudEdnYUN6eHpQWDhUa2tHWkdOT0hQNDNQVmVwYnhtL0llNWQwcFM2ZnJRRnBKbkYzd3kiLCJleHAiOjE3MjY2Njk4MzR9.OJMvLTXG7oYJ8T4aw0_gPVGAUSqowKu8lARoaDNHU2o77ehY9na2jh123tenRduDiu9fLP7ANLCQm_XHnVu8RA",
// 		"refresh_token": "Jy3SU-ToorSu-UPfKleKov4C2mBsiKFnjygsysTpYcg=",
// 	}
// 	jsonBody, _ := json.Marshal(requestBody)
// 	req := httptest.NewRequest(http.MethodPost, "/token/refresh", bytes.NewReader(jsonBody))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()

// 	// Execute request
// 	e.ServeHTTP(rec, req)

// 	// Assert results
// 	assert.Equal(t, http.StatusOK, rec.Code)

// 	var tokenPair token.TokenPair
// 	err := json.Unmarshal(rec.Body.Bytes(), &tokenPair)
// 	assert.NoError(t, err)

// 	assert.Equal(t, "newMockAccessToken", tokenPair.AccessToken)
// 	assert.Equal(t, "newMockRefreshToken", tokenPair.RefreshToken)
// }
