package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerListChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Error fetching all Chirps: %s", err)
		respondWithInternalServerError(w)
		return
	}

	payload := []Chirp{}

	for _, chirp := range chirps {
		nextChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		}
		payload = append(payload, nextChirp)
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
