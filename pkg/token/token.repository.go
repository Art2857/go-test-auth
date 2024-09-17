package token

import (
	"auth-service/pkg/database"
	"context"
	"errors"
)

func RefreshTokenCreate(refreshTokenHash, ip string) (bool, error) {
	_, err := database.DB.Exec(context.Background(), "INSERT INTO refresh_tokens (token_hash, ip) VALUES ($1, $2)", refreshTokenHash, ip)
	if err != nil {
		return false, errors.New("Refresh Token Create Database error: " + err.Error())
	}

	return true, nil
}

func RefreshTokenGetIP(refreshTokenHash string) (string, error) {
	var ip string

	err := database.DB.QueryRow(context.Background(), "SELECT ip FROM refresh_tokens WHERE token_hash = $1", refreshTokenHash).Scan(&ip)
	if err != nil {
		return "", errors.New("Refresh Token Get Database error: " + err.Error())
	}

	return ip, nil
}

func RefreshTokenRemove(refreshTokenHash string) (bool, error) {
	var ip string

	err := database.DB.QueryRow(context.Background(), "SELECT ip FROM refresh_tokens WHERE token_hash = $1", refreshTokenHash).Scan(&ip)
	if err != nil {
		return false, errors.New("Refresh Token Remove Database error: " + err.Error())
	}

	_, err = database.DB.Exec(context.Background(), "DELETE FROM refresh_tokens WHERE token_hash = $1", refreshTokenHash)
	if err != nil {
		return false, errors.New("Refresh Token Remove Database error: " + err.Error())
	}

	return true, nil
}
