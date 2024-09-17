package token

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"auth-service/pkg/config"
	"auth-service/pkg/mail"
)

var jwtSecret = []byte(config.Env.JWT_SECRET)

// Структура JWT токена
type Claims struct {
	UserID           string `json:"user_id"`
	IP               string `json:"ip"`
	RefreshTokenHash string `json:"refresh_token_hash"` // Необходим для связки с access токеном

	jwt.StandardClaims
}

// Функция для генерации Access токена (JWT)
func signAccessToken(userID, ip, refreshToken string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	// Хешируем refresh токен для связки с access токеном
	refreshTokenHash, err := hashRefreshTokenToJWT(refreshToken)
	if err != nil {
		log.Print(err)
		return "", err
	}

	claims := &Claims{
		UserID:           userID,
		IP:               ip,
		RefreshTokenHash: refreshTokenHash,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	return token.SignedString(jwtSecret)
}

// Функция для верификации Access токена
func VerifyAccessToken(token string, ignoreExpiration bool) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		log.Print(err)
		return nil, err
	}

	if claims, ok := tokenClaims.Claims.(*Claims); ok && (tokenClaims.Valid || ignoreExpiration) {
		return claims, nil
	}

	return nil, err
}

// Функция для генерации случайного refresh токена в формате base64
func generateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

// Хеширование refresh токена для связки с access токеном
func hashRefreshTokenToJWT(refreshToken string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return string(hash), nil
}

// Хеширование refresh токена для хранения в базе данных
func hashRefreshTokenToDatabase(refreshToken string) string {
	hash := sha256.Sum256([]byte(refreshToken + "refresh token database salt"))
	return hex.EncodeToString(hash[:])
}

// Функция для генерации пары Access и Refresh токенов
func GenerateTokenPair(userID, ip string) (map[string]string, error) {
	refreshToken, err := generateRefreshToken() // Генерируем refresh token
	if err != nil {
		log.Print(err)
		return nil, err
	}

	accessToken, err := signAccessToken(userID, ip, refreshToken) // Генерируем access token с идентификатором пользователя и ip-адресом
	if err != nil {
		log.Print(err)
		return nil, errors.New("Error generating access token: " + err.Error())
	}

	refreshTokenHash := hashRefreshTokenToDatabase(refreshToken) // Хешируем refresh token

	// Сохранение refresh токена в базу
	_, err = RefreshTokenCreate(refreshTokenHash, ip)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Отправляем токены клиенту
	tokenPair := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	return tokenPair, nil
}

// Функция для обновления пары Access и Refresh токенов
func RefreshTokenPair(accessToken, refreshToken, ip string) (map[string]string, error) {
	claims, err := VerifyAccessToken(accessToken, true)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Сверяем что refresh связан с access токеном
	err = bcrypt.CompareHashAndPassword([]byte(claims.RefreshTokenHash), []byte(refreshToken))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	refreshTokenHashDatabase := hashRefreshTokenToDatabase(refreshToken)

	// Подбираем ip, который был при аутентификации
	beforeIP, err := RefreshTokenGetIP(refreshTokenHashDatabase)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Если IP-адрес пользователя изменился, то отправляем предупреждение на почту (моковые данные)
	if ip != beforeIP {
		log.Print(err)
		mail.SendEmail("mock@ya.ru", "Invalid IP", "Warning: Invalid IP")

		return nil, errors.New("invalid IP")
	}

	_, err = RefreshTokenRemove(refreshTokenHashDatabase)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	tokenPair, err := GenerateTokenPair(claims.UserID, ip)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return tokenPair, err
}
