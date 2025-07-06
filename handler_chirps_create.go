package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rdcassin/server-go/internal/database"

	"github.com/google/uuid"
)

const charLimit int = 140

var profaneWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerAddChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type returnVals struct {
		Chirp
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params in handlerAddChirp: %s", err)
		respondWithInternalServerError(w)
		return
	}

	params.Body, err = cleanChirp(params.Body)
	if err != nil {
		log.Printf("Error validating Chirp: %s", err)
		msg := fmt.Sprintf("%s", err)
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	newChirpParams := database.CreateChirpParams{
		Body: params.Body,
		UserID: params.UserID,
	}

	newChirp, err := cfg.db.CreateChirp(r.Context(), newChirpParams)
	if err != nil {
		log.Printf("Error creating new Chirp: %s", err)
		respondWithInternalServerError(w)
		return
	}

	payload := returnVals{
		Chirp: Chirp{
			ID: newChirp.ID,
			CreatedAt: newChirp.CreatedAt,
			UpdatedAt: newChirp.UpdatedAt,
			Body: newChirp.Body,
			UserID: newChirp.UserID,
		},
	}
	
	respondWithJSON(w, http.StatusCreated, payload)
}

func validateChirp(paramsBody string) (string, error) {
	chars := []rune(paramsBody)
	if len(chars) > charLimit {
		return "", fmt.Errorf("Chirp exceeds maximum character limit")
	}

	return paramsBody, nil
}

func cleanChirp(paramsBody string) (string, error) {
	words := strings.Split(paramsBody, " ")

	cleanBody := []string{}

	for _, word := range words {
		if _, exists := profaneWords[strings.ToLower(word)]; exists {
			cleanBody = append(cleanBody, "****")
		} else {
			cleanBody = append(cleanBody, word)
		}
	}

	return validateChirp(strings.Join(cleanBody, " "))
}