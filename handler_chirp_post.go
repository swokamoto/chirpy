package main

import (
	"encoding/json"
	"net/http"
	"time"
	"strings"
	"log"
	"github.com/google/uuid"
	"github.com/swokamoto/chirpy/internal/database"
)

func (cfg *apiConfig) handlerChirpCreate(w http.ResponseWriter, r *http.Request) {
	log.Println("handlerChirpCreate called")
	
	type parameters struct {
		Body string    `json:"body"`
		UserID  uuid.UUID `json:"user_id"`
	}
	type response struct {
		ChirpID uuid.UUID `json:"chirp_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body    string    `json:"body"`
		UserID	uuid.UUID `json:"user_id"`
		
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	log.Printf("Decoded params: %+v", params)

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	cleaned := getCleanedBody(params.Body, badWords)
	log.Printf("Cleaned body: %s", cleaned)

	createParams := database.CreateChirpParams{
        Body: cleaned,
        UserID: uuid.NullUUID{
            UUID:  params.UserID,
            Valid: true,
        },
    }
    log.Printf("About to call CreateChirp with params: %+v", createParams)
    

	chirp, err := cfg.db.CreateChirp(r.Context(), createParams)
    if err != nil {
        log.Printf("CreateChirp error: %v", err)
        respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
        return
    }
    
    log.Printf("Successfully created chirp: %+v", chirp)
	

	userID := uuid.UUID{}
    if chirp.UserID.Valid {
        userID = chirp.UserID.UUID
    }

	respondWithJSON(w, http.StatusCreated, response{
		ChirpID:   chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      cleaned,
		UserID:    userID,
		
	})	
}

func getCleanedBody(body string, badWords map[string]struct{}) string {
	words := strings.Split(body, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}