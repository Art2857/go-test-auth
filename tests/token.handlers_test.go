package token_test

import (
	"auth/pkg/database/models"
	"auth/pkg/token"
	test_utils "auth/tests/utils"
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"bytes"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Написал часть тестов, так как много чего ещё дописывать нужно, чтобы всё затестить...

// Кейсы для тестирования:
// Должен вернуть код 200
// Должен вернуть access_token и refresh_token
// Должен вернуть access_token типа jwt и refresh_token типа base64
// Должен вернуть claims(jwt) валидируя пару access и refresh токенов как связанные
// Должен валидировать, что claims(jwt) содержит: userID и является uuid, RefreshTokenHash - является строкой
// Запись должна появиться в базе данных
func TestHandlerGenerateTokenPair(t *testing.T) {
	app := test_utils.SetupTest()
	defer test_utils.EndTest(app)

	// Test request data
	requestBody := map[string]string{
		"user_id": "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
	}

	// Execute request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/token/login", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	app.Router.ServeHTTP(rec, req)

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

	_, err = base64.URLEncoding.DecodeString(tokenPair.RefreshToken)
	assert.NoError(t, err, "Refresh token should be a valid Base64 encoded string")

	claims, err := app.TokenService.VerifyTokenPair(tokenPair)
	assert.NoError(t, err, "TokenPair verification should not return an error")
	assert.NotNil(t, claims, "Claims should not be nil")

	assert.NotEmpty(t, claims.UserID, "UserID should not be empty")
	assert.NotEmpty(t, claims.RefreshTokenHash, "RefreshTokenHash should not be empty")

	_, err = uuid.Parse(claims.UserID)
	assert.NoError(t, err, "UserID should be a valid UUID")
	assert.IsType(t, "", claims.RefreshTokenHash, "RefreshTokenHash should be of type string")

	var refreshTokenModel models.RefreshTokenModel
	if err := app.DB.Postgres.DB.First(&refreshTokenModel, "token_hash = ?", app.TokenService.HashRefreshTokenToDatabase(refreshToken)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			assert.Fail(t, "Refresh token should be in the database", err)
		} else {
			assert.Fail(t, "Database error", err)
		}
	}
}

// При верных данных идёт проверка из первого теста
func TestHandlerRefreshTokenPair(t *testing.T) {
	// Setup
	app := test_utils.SetupTest()
	defer test_utils.EndTest(app)

	ip := "192.0.2.1"
	mockRefreshToken := "5Zc5pf9ureStwu_6yg2CwcuFERIivaXnlLGJfX8XMwA="

	mockAccessToken, err := app.TokenService.SignAccessToken("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11", ip, mockRefreshToken)
	if err != nil {
		assert.Fail(t, "Setup error - failed to create access token", err)
	}

	refreshTokenHash := app.TokenService.HashRefreshTokenToDatabase(mockRefreshToken)

	_, err = app.TokenRepository.RefreshTokenCreate(refreshTokenHash, ip)
	if err != nil {
		assert.Fail(t, "Setup error - failed to create refresh token", err)
	}

	// Test request data

	requestBody := map[string]string{
		"access_token":  mockAccessToken,
		"refresh_token": mockRefreshToken,
	}

	fmt.Println(requestBody)

	// Execute request
	jsonBody, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/token/refresh", bytes.NewReader(jsonBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	app.Router.ServeHTTP(rec, req)

	// Assert results (Копия проверки из первого теста)
	assert.Equal(t, http.StatusOK, rec.Code)

	var tokenPair *token.TokenPair
	err = json.Unmarshal(rec.Body.Bytes(), &tokenPair)
	assert.NoError(t, err)

	accessToken, refreshToken := tokenPair.AccessToken, tokenPair.RefreshToken

	assert.NotEmpty(t, accessToken, "Access token should not be empty")
	assert.NotEmpty(t, refreshToken, "Refresh token should not be empty")

	_, _, err = new(jwt.Parser).ParseUnverified(tokenPair.AccessToken, jwt.MapClaims{})
	assert.NoError(t, err, "Access token should be a valid JWT")

	_, err = base64.URLEncoding.DecodeString(tokenPair.RefreshToken)
	assert.NoError(t, err, "Refresh token should be a valid Base64 encoded string")

	claims, err := app.TokenService.VerifyTokenPair(tokenPair)
	assert.NoError(t, err, "TokenPair verification should not return an error")
	assert.NotNil(t, claims, "Claims should not be nil")

	assert.NotEmpty(t, claims.UserID, "UserID should not be empty")
	assert.NotEmpty(t, claims.RefreshTokenHash, "RefreshTokenHash should not be empty")

	_, err = uuid.Parse(claims.UserID)
	assert.NoError(t, err, "UserID should be a valid UUID")
	assert.IsType(t, "", claims.RefreshTokenHash, "RefreshTokenHash should be of type string")

	var refreshTokenModel models.RefreshTokenModel
	if err := app.DB.Postgres.DB.First(&refreshTokenModel, "token_hash = ?", app.TokenService.HashRefreshTokenToDatabase(refreshToken)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			assert.Fail(t, "Refresh token should be in the database", err)
		} else {
			assert.Fail(t, "Database error", err)
		}
	}
}

// Также стоит рассмотреть такие кейсы:
// При другом ip должен быть отправлено предупреждение на почту
// При неверном access ошибка, при неверном refresh ошибка
