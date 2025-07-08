package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerListChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error fetching all Chirps: %s", err)
		respondWithInternalServerError(w)
		return
	}

	// Grab optional author_id query
	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		authorID, err = uuid.Parse(authorIDString)
		if err != nil {
			log.Printf("Error parsing ID: %s", err)
			respondWithInternalServerError(w)
			return
		}
	}

	payload := []Chirp{}
	for _, chirp := range dbChirps {
		if authorID != uuid.Nil && chirp.UserID != authorID {
			continue
		}
		newChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		payload = append(payload, newChirp)
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func (cfg *apiConfig) handlerFetchChirp(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Chirp
	}

	chirp, err := cfg.getChirp(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	payload := returnVals{
		Chirp: Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		},
	}

	respondWithJSON(w, http.StatusOK, payload)
}
