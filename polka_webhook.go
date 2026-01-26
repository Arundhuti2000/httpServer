package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Arundhuti2000/httpserver/internal/auth"
	"github.com/google/uuid"
)

type PolkaWebhookRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	apiKey, err := GetAPIKeyFromHeaders(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req PolkaWebhookRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Only process user.upgraded events
	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Upgrade the user to Chirpy Red
	_, err = cfg.DB.UpgradeUserToChirpyRed(r.Context(), req.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to upgrade user")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
// GetAPIKeyFromHeaders extracts the API key from the Authorization header
func GetAPIKeyFromHeaders(headers http.Header) (string, error) {
	return auth.GetAPIKey(headers)
}