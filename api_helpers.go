package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rdcassin/server-go/internal/database"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, out interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(out); err != nil {
		log.Printf("Error decoding JSON body: %s", err)
		respondWithInternalServerError(w)
		return false
	}
	return true
}

func (cfg *apiConfig) getChirp(r *http.Request) (database.Chirp, error) {
	rawChirpID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(rawChirpID)
	if err != nil {
		log.Printf("Error parsing Chirp ID %s: %s", rawChirpID, err)
		return database.Chirp{}, err
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Chirp with ID %s does not exist", chirpID)
			return database.Chirp{}, err
		}
		log.Printf("Error fetching Chirp in handlerFetchChirp with ID: %s: %s", chirpID, err)
		return database.Chirp{}, err
	}

	return chirp, nil
}