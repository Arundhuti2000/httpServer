package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Arundhuti2000/httpserver/internal/auth"
	"github.com/Arundhuti2000/httpserver/internal/database"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}


func cleanProfanity(s string) string {
	banned := []string{"kerfuffle", "sharbert", "fornax"}
	parts := strings.Split(s, " ")
	for i, tok := range parts {
		for _, b := range banned {
			if strings.EqualFold(tok, b) {
				parts[i] = "****"
				break
			}
		}
	}
	return strings.Join(parts, " ")
}





func (cfg *apiConfig) handlerusers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Hash the password
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}
	
	// Create user in database
	dbUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create user")
		return
	}

	// Convert database.User to main.User for JSON response
	user := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email             string `json:"email"`
		Password          string `json:"password"`
		ExpiresInSeconds  *int   `json:"expires_in_seconds"`
	}

	type response struct {
		ID        string `json:"id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		Email     string `json:"email"`
		Token     string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Look up user by email
	dbUser, err := cfg.DB.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error getting user: %s", err)
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// Check password
	match, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil || !match {
		log.Printf("Error checking password or password mismatch")
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// Determine expiration time
	expiresIn := time.Hour // Default: 1 hour
	if params.ExpiresInSeconds != nil {
		expiresIn = time.Duration(*params.ExpiresInSeconds) * time.Second
		// Cap at 1 hour
		if expiresIn > time.Hour {
			expiresIn = time.Hour
		}
	}

	// Create JWT token
	token, err := auth.MakeJWT(dbUser.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		log.Printf("Error creating JWT: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create token")
		return
	}

	// Return user with token
	resp := response{
		ID:        dbUser.ID.String(),
		CreatedAt: dbUser.CreatedAt.Format(time.RFC3339),
		UpdatedAt: dbUser.UpdatedAt.Format(time.RFC3339),
		Email:     dbUser.Email,
		Token:     token,
	}

	respondWithJSON(w, http.StatusOK, resp)
}

