package token

import (
	"auth-service/pkg/utils"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TokenHandlers struct {
	TokenService *TokenService
}

func NewHanders(tokenService *TokenService) *TokenHandlers {
	return &TokenHandlers{TokenService: tokenService}
}

func (h *TokenHandlers) InitHandlers(e *echo.Echo) *echo.Group {

	tokenGroup := e.Group("/token")
	{
		tokenGroup.POST("/login", h.HandlerGenerateTokenPair)
		tokenGroup.POST("/refresh", h.HandlerRefreshTokenPair)
	}

	return tokenGroup
}

// HandlerGenerateTokenPair генерирует пару токенов (access и refresh)
// @Summary      Генерация пары токенов
// @Description  Генерирует access и refresh токены для пользователя по его ID.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  object{user_id=string}  true  "Тело запроса с user_id"
// @Success      200  {object}  TokenPair  "Пара токенов"
// @Failure      400  {string}  string  "Неверный запрос"
// @Failure      500  {string}  string  "Ошибка генерации токенов"
// @Router       /token/login [post]
func (h *TokenHandlers) HandlerGenerateTokenPair(c echo.Context) error {
	ip := c.RealIP()

	var data struct {
		UserID string `json:"user_id"`
	}

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	if !utils.IsValidUUID(data.UserID) {
		return c.JSON(http.StatusBadRequest, "Invalid UUID format")
	}

	tokenPair, err := h.TokenService.GenerateTokenPair(data.UserID, ip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error generating token pair")
	}

	return c.JSON(http.StatusOK, tokenPair)
}

// HandlerRefreshTokenPair обновляет access и refresh токены
// @Summary      Обновление пары токенов
// @Description  Обновляет access и refresh токены, используя действующую refresh-токен
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body  object{access_token=string,refresh_token=string}  true  "Тело запроса с токенами"
// @Success      200  {object}  TokenPair  "Обновленная пара токенов"
// @Failure      400  {string}  string  "Неверный запрос"
// @Failure      500  {string}  string  "Ошибка обновления токенов"
// @Router       /token/refresh [post]
func (h *TokenHandlers) HandlerRefreshTokenPair(c echo.Context) error {
	ip := c.RealIP()

	var data struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.Bind(&data); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid request body")
	}

	tokenPair, err := h.TokenService.RefreshTokenPair(&TokenPair{AccessToken: data.AccessToken, RefreshToken: data.RefreshToken}, ip)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error refreshing token pair")
	}

	return c.JSON(http.StatusOK, tokenPair)
}
