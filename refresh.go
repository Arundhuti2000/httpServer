package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Arundhuti2000/httpserver/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	// Extract refresh token from Authorization header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "invalid or missing token")
		return
	}

	// Look up user from refresh token (validates token exists, not expired, not revoked)
	dbUser, err := cfg.DB.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Error getting user from refresh token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "invalid or expired token")
		return
	}

	// Create new access token that expires in 1 hour
	accessToken, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create token")
		return
	}

	// Return new access token
	resp := response{
		Token: accessToken,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Extract refresh token from Authorization header
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "invalid or missing token")
		return
	}

	// Revoke the token in the database
	err = cfg.DB.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		log.Printf("Error revoking refresh token: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke token")
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
