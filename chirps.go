package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerchirps(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
        Body string `json:"body"`
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

