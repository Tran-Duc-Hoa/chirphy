package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/Tran-Duc-Hoa/chirphy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token        string `json:"token"`
	}
	// Get the Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusUnauthorized, "Authorization header is required", nil)
		return
	}

	// Check if the token is a Bearer token
	if !strings.HasPrefix(authHeader, "Bearer ") {
		respondWithError(w, http.StatusUnauthorized, "Invalid token format", nil)
		return
	}

	// Extract the refresh token
	refreshToken := strings.TrimPrefix(authHeader, "Bearer ")
	if refreshToken == "" {
		respondWithError(w, http.StatusUnauthorized, "Invalid refresh token", nil)
		return
	}

	// Validate and process the refresh token
	token, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get refresh token", err)
		return
	}
	
	// Check if the token exists and is not expired
	if token.RevokedAt.Valid || token.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "Invalid or expired refresh token", err)
		return
	}

	accessToken, err := auth.MakeJWT(
		token.UserID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}
	

	respondWithJSON(w, http.StatusOK, response{
		Token: accessToken,
	})
}