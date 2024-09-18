package token

import (
	"auth/pkg/database/models"
	"errors"

	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{
		db: db,
	}
}

func (repo *TokenRepository) RefreshTokenCreate(refreshTokenHash, ip string) (bool, error) {
	refreshToken := models.RefreshTokenModel{
		TokenHash: refreshTokenHash,
		IP:        ip,
	}

	if err := repo.db.Create(&refreshToken).Error; err != nil {
		return false, errors.New("Refresh Token Create Database error: " + err.Error())
	}

	return true, nil
}

func (repo *TokenRepository) RefreshTokenGetIP(refreshTokenHash string) (string, error) {
	var refreshToken models.RefreshTokenModel

	if err := repo.db.Where("token_hash = ?", refreshTokenHash).First(&refreshToken).Error; err != nil {
		return "", errors.New("Refresh Token Get Database error: " + err.Error())
	}

	return refreshToken.IP, nil
}

func (repo *TokenRepository) RefreshTokenRemove(refreshTokenHash string) (bool, error) {
	result := repo.db.Where("token_hash = ?", refreshTokenHash).Unscoped().Delete(&models.RefreshTokenModel{})

	if result.Error != nil {
		return false, errors.New("Refresh Token Remove Database error: " + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
