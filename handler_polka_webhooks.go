package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	type webhookData struct {
		UserID string `json:"user_id"`
	}

	type webhookRequest struct {
		Event string      `json:"event"`
		Data  webhookData `json:"data"`
	}

	// Parse the webhook request
	decoder := json.NewDecoder(r.Body)
	webhookReq := webhookRequest{}
	err := decoder.Decode(&webhookReq)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode webhook request", err)
		return
	}

	// Only handle user.upgraded events, ignore all others
	if webhookReq.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Parse the user ID
	userID, err := uuid.Parse(webhookReq.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	// Upgrade the user to Chirpy Red
	_, err = cfg.db.MakeUserRed(r.Context(), userID)
	if err != nil {
		// User not found or other database error
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	// Success - return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
