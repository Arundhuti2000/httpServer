package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Arundhuti2000/httpserver/internal/database"
	"github.com/google/uuid"
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

func (cfg *apiConfig) handlerchirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
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

	// Create chirp in database
	dbChirp, err := cfg.DB.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned,
		UserID: params.UserID,
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



func (cfg *apiConfig) handlerusers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
        Email string `json:"email"`
    }
	
	decoder := json.NewDecoder(r.Body)
	params:=parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, http.StatusBadRequest, "invalid request body")
        return
    } 
	
	// Create user in database
	dbUser, err := cfg.DB.CreateUser(r.Context(), params.Email)
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

