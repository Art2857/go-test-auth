package token

import (
	"auth/pkg/database"
	"auth/pkg/database/models"
	"fmt"
)

type TokenRepository struct {
	db *database.Postgres
}

func NewRepository(db *database.Postgres) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

// Ошибка создания Refresh Token в базе данных
type ErrRefreshTokenCreate struct {
	Reason string
}

func (e *ErrRefreshTokenCreate) Error() string {
	return fmt.Sprintf("Refresh Token creation error in the database: %s", e.Reason)
}

func (repo *TokenRepository) RefreshTokenCreate(refreshTokenHash, ip string) (bool, error) {
	refreshToken := models.RefreshTokenModel{
		TokenHash: refreshTokenHash,
		IP:        ip,
	}

	if err := repo.db.Create(&refreshToken).Error; err != nil {
		return false, &ErrRefreshTokenCreate{Reason: err.Error()}
	}

	return true, nil
}

// Ошибка получения IP по Refresh Token
type ErrRefreshTokenGetIP struct {
	Reason string
}

func (e *ErrRefreshTokenGetIP) Error() string {
	return fmt.Sprintf("Error retrieving IP by Refresh Token: %s", e.Reason)
}

func (repo *TokenRepository) RefreshTokenGetIP(refreshTokenHash string) (string, error) {
	var refreshToken models.RefreshTokenModel

	if err := repo.db.Where("token_hash = ?", refreshTokenHash).First(&refreshToken).Error; err != nil {
		return "", &ErrRefreshTokenGetIP{Reason: err.Error()}
	}

	return refreshToken.IP, nil
}

// Ошибка удаления Refresh Token из базы данных
type ErrRefreshTokenRemove struct {
	Reason string
}

func (e *ErrRefreshTokenRemove) Error() string {
	return fmt.Sprintf("Error deleting Refresh Token from the database: %s", e.Reason)
}

func (repo *TokenRepository) RefreshTokenRemove(refreshTokenHash string) (bool, error) {
	result := repo.db.Where("token_hash = ?", refreshTokenHash).Unscoped().Delete(&models.RefreshTokenModel{})

	if result.Error != nil {
		return false, &ErrRefreshTokenRemove{Reason: result.Error.Error()}
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
