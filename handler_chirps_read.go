package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerListChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error fetching all Chirps: %s", err)
		msg := "Error fetching all Chirps"
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	payload := []Chirp{}

	for _, chirp := range chirps {
		nextChirp := Chirp {
			ID: chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body: chirp.Body,
			UserID: chirp.UserID,
		}
		payload = append(payload, nextChirp)
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func (cfg *apiConfig) handlerFetchChirp(w http.ResponseWriter, r *http.Request) {	
	rawChirpID := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(rawChirpID)
	if err != nil {
		log.Printf("Error parsing Chirp ID %s: %s", rawChirpID, err)
		msg := "Error parsing Chirp ID"
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Chirp with ID %s does not exist", chirpID)
			msg := fmt.Sprintf("Chirp with ID %s does not exist", chirpID)
			respondWithError(w, http.StatusNotFound, msg)
			return
		}
		log.Printf("Error fetching Chirp with ID: %s: %s", chirpID, err)
		msg := fmt.Sprintf("Error fetching Chirp with ID: %v", chirpID)
		respondWithError(w, http.StatusInternalServerError, msg)
		return 
	}

	payload := Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}

	respondWithJSON(w, http.StatusOK, payload)
}