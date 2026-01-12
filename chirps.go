package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Arundhuti2000/httpserver/internal/auth"
	"github.com/Arundhuti2000/httpserver/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerchirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	
	// Extract token from Authorization header
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	// Validate JWT token
	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		log.Printf("Error validating JWT: %s", err)
		respondWithError(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	
	// Validate chirp length
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	// Clean profanity
	cleaned := cleanProfanity(params.Body)

	// Create chirp in database with the authenticated user's ID
	dbChirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: userID,
	})
	if err != nil {
		log.Printf("Error creating chirp: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't create chirp")
		return
	}

	// Convert database.Chirp to main.Chirp for JSON response
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}

func (cfg *apiConfig) handlerGetAllChirps(w http.ResponseWriter, r *http.Request) {
	// Get all chirps from database
	dbChirps, err := cfg.DB.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error retrieving chirps: %s", err)
		respondWithError(w, http.StatusInternalServerError, "couldn't retrieve chirps")
		return
	}

	// Convert []database.Chirp to []Chirp for JSON response
	chirps := make([]Chirp, len(dbChirps))
	for i, dbChirp := range dbChirps {
		chirps[i] = Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserID:    dbChirp.UserID,
		}
	}

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpByID(w http.ResponseWriter, r *http.Request) {
	// Get chirp ID from path parameter
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid chirp ID")
		return
	}

	// Get chirp from database
	dbChirp, err := cfg.DB.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "chirp not found")
		return
	}

	// Convert database.Chirp to main.Chirp for JSON response
	chirp := Chirp{
		ID:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserID:    dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, chirp)
}