package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Arundhuti2000/httpserver/internal/auth"
	"github.com/Arundhuti2000/httpserver/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Extract token from Authorization header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Validate JWT token to get user ID
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT: %s", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Parse request body
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Hash the new password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't update user")
		return
	}

	// Update user in database
	dbUser, err := cfg.DB.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		log.Printf("Error updating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't update user")
		return
	}

	// Convert database.User to main.User for JSON response (omitting password)
	user := User{
		ID:          dbUser.ID,
		CreatedAt:   dbUser.CreatedAt,
		UpdatedAt:   dbUser.UpdatedAt,
		Email:       dbUser.Email,
		IsChirpyRed: dbUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, user)
}
