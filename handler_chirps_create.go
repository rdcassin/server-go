package main

import (
	"fmt"
	"github.com/rdcassin/server-go/internal/database"
	"log"
	"net/http"
	"strings"
)

const charLimit int = 140

var profaneWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func (cfg *apiConfig) handlerAddChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		Chirp
	}

	userID, err := cfg.validateUser(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := parameters{}
	if !decodeJSONBody(w, r, &params) {
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
		Body:   params.Body,
		UserID: userID,
	}

	newChirp, err := cfg.db.CreateChirp(r.Context(), newChirpParams)
	if err != nil {
		log.Printf("Error creating new Chirp: %s", err)
		respondWithInternalServerError(w)
		return
	}

	payload := returnVals{
		Chirp: Chirp{
			ID:        newChirp.ID,
			CreatedAt: newChirp.CreatedAt,
			UpdatedAt: newChirp.UpdatedAt,
			Body:      newChirp.Body,
			UserID:    newChirp.UserID,
		},
	}

	respondWithJSON(w, http.StatusCreated, payload)
}

func validateChirp(paramsBody string) (string, error) {
	chars := []rune(paramsBody)
	if len(chars) > charLimit {
		return "", fmt.Errorf("chirp exceeds maximum character limit")
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
