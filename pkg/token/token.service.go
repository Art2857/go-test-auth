package token

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	"auth/pkg/config"
	"auth/pkg/mail"
	"auth/pkg/utils"
)

type TokenService struct {
	ConfigService   *config.Config
	TokenRepository *TokenRepository
	MailService     *mail.MailService
}

func NewService(configService *config.Config, tokenRepository *TokenRepository, mailService *mail.MailService) *TokenService {
	return &TokenService{
		ConfigService:   configService,
		TokenRepository: tokenRepository,
		MailService:     mailService,
	}
}

// TokenPair содержит пару токенов
// @Description Структура для пару токенов access и refresh
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// var jwtSecret = []byte(config.Env.JWTSecret)

// Структура JWT токена
type Claims struct {
	UserID           string `json:"user_id"`
	IP               string `json:"ip"`
	RefreshTokenHash string `json:"refresh_token_hash"` // Необходим для связки с access токеном

	jwt.StandardClaims
}

// Функция для генерации Access токена (JWT)
func (s *TokenService) SignAccessToken(userID, ip, refreshToken string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	// Хешируем refresh токен для связки с access токеном
	refreshTokenHash, err := s.hashRefreshTokenToJWT(refreshToken)
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

	return token.SignedString([]byte(s.ConfigService.Env.JWTSecret))
}

// Функция для верификации Access токена
func (s *TokenService) VerifyAccessToken(token string, ignoreExpiration bool) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.ConfigService.Env.JWTSecret), nil
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

func (s *TokenService) VerifyTokenPair(TokenPair *TokenPair) (*Claims, error) {
	claims, err := s.VerifyAccessToken(TokenPair.AccessToken, true)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Сверяем что refresh связан с access токеном
	err = bcrypt.CompareHashAndPassword([]byte(claims.RefreshTokenHash), []byte(TokenPair.RefreshToken))
	if err != nil {
		log.Print(err)
		return nil, err
	}

	return claims, nil
}

// Функция для генерации случайного refresh токена в формате base64
func (s *TokenService) GenerateRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)

	_, err := rand.Read(tokenBytes)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}

// Хеширование refresh токена для связки с access токеном
func (s *TokenService) hashRefreshTokenToJWT(refreshToken string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return string(hash), nil
}

// Хеширование refresh токена для хранения в базе данных
func (s *TokenService) HashRefreshTokenToDatabase(refreshToken string) string {
	return utils.Sum256(refreshToken + "refresh token database salt")
}

// Функция для генерации пары Access и Refresh токенов
func (s *TokenService) GenerateTokenPair(userID, ip string) (*TokenPair, error) {
	refreshToken, err := s.GenerateRefreshToken() // Генерируем refresh token
	if err != nil {
		log.Print(err)
		return nil, err
	}

	accessToken, err := s.SignAccessToken(userID, ip, refreshToken) // Генерируем access token с идентификатором пользователя и ip-адресом
	if err != nil {
		log.Print(err)
		return nil, errors.New("Error generating access token: " + err.Error())
	}

	refreshTokenHash := s.HashRefreshTokenToDatabase(refreshToken) // Хешируем refresh token

	// Сохранение refresh токена в базу
	_, err = s.TokenRepository.RefreshTokenCreate(refreshTokenHash, ip)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	// Отправляем токены клиенту
	tokenPair := &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return tokenPair, nil
}

// Функция для обновления пары Access и Refresh токенов
func (s *TokenService) RefreshTokenPair(tokenPair *TokenPair, ip string) (*TokenPair, error) {
	claims, err := s.VerifyTokenPair(tokenPair)
	if err != nil {
		log.Print("Error verifying token pair: ", err)
		return nil, err
	}

	refreshTokenHashDatabase := s.HashRefreshTokenToDatabase(tokenPair.RefreshToken)

	// Подбираем ip, который был при аутентификации
	beforeIP, err := s.TokenRepository.RefreshTokenGetIP(refreshTokenHashDatabase)
	if err != nil {
		log.Print("Error getting IP: ", err)
		return nil, err
	}

	// Если IP-адрес пользователя изменился, то отправляем предупреждение на почту (моковые данные)
	if ip != beforeIP {
		log.Print("Warning: IP changed from " + beforeIP + " to " + ip)
		s.MailService.SendEmail("mock@ya.ru", "Invalid IP", "Warning: Invalid IP")

		return nil, errors.New("Warning: IP changed from " + beforeIP + " to " + ip)
	}

	exists, err := s.TokenRepository.RefreshTokenRemove(refreshTokenHashDatabase)
	if err != nil {
		log.Print("Error removing refresh token: ", err)
		return nil, err
	}

	if !exists {
		log.Print("Refresh token not found")
		return nil, errors.New("refresh token not found")
	}

	newTokenPair, err := s.GenerateTokenPair(claims.UserID, ip)
	if err != nil {
		log.Print("Error generating token pair: ", err)
		return nil, err
	}

	return newTokenPair, err
}
