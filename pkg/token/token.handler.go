package token

import (
	"encoding/json"
	"net/http"
)

func HandlerGenerateTokenPair(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr

	var data struct {
		UserID string `json:"user_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokenPair, err := GenerateTokenPair(data.UserID, ip)
	if err != nil {
		http.Error(w, "Error generating token pair", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tokenPair)
}

func HandlerRefreshTokenPair(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr

	var data struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	tokenPair, err := RefreshTokenPair(data.AccessToken, data.RefreshToken, ip)
	if err != nil {
		http.Error(w, "Error refreshing token pair", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tokenPair)
}
