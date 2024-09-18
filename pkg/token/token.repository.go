package token

import (
	"auth-service/pkg/database"
	"errors"
)

func RefreshTokenCreate(refreshTokenHash, ip string) (bool, error) {
	refreshToken := database.RefreshTokenModel{
		TokenHash: refreshTokenHash,
		IP:        ip,
	}

	if err := database.DB.Create(&refreshToken).Error; err != nil {
		return false, errors.New("Refresh Token Create Database error: " + err.Error())
	}

	return true, nil
}

func RefreshTokenGetIP(refreshTokenHash string) (string, error) {
	var refreshToken database.RefreshTokenModel

	if err := database.DB.Where("token_hash = ?", refreshTokenHash).First(&refreshToken).Error; err != nil {
		return "", errors.New("Refresh Token Get Database error: " + err.Error())
	}

	return refreshToken.IP, nil
}

func RefreshTokenRemove(refreshTokenHash string) (bool, error) {
	result := database.DB.Where("token_hash = ?", refreshTokenHash).Unscoped().Delete(&database.RefreshTokenModel{})

	if result.Error != nil {
		return false, errors.New("Refresh Token Remove Database error: " + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return false, nil
	}

	return true, nil
}
