package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func writeJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
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
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
        return
		
    } 
	if len(params.Body) > 140 {
        writeJSONError(w, http.StatusBadRequest, "Chirp is too long")
        return
    }
	type returnVals struct {
		Valid bool `json:"valid"`
	}

	respBody := returnVals{Valid: true}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(respBody); err != nil {
        log.Printf("Error encoding response: %v", err)
        writeJSONError(w, http.StatusInternalServerError, "Something went wrong")
        return
    }
}