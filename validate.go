package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

func (cfg *apiConfig) handlerValidateChirps(w http.ResponseWriter, r *http.Request) {
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
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned := cleanProfanity(params.Body)

	resp := map[string]string{"cleaned_body": cleaned}
	respondWithJSON(w, http.StatusOK, resp)
}