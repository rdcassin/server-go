package main

import (
	"encoding/json"
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

func handlerChirp(w http.ResponseWriter, r *http.Request) {
	return
}

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Fatalf("Error decoding params in handlerChirpValidate: %s", err)
		msg := "Something went wrong"
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	params.Body = cleanChirp(params.Body)

	chars := []rune(params.Body)
	if len(chars) > charLimit {
		msg := "Chirp is too long"
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	payload := returnVals{
		CleanedBody: params.Body,
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func cleanChirp(paramsBody string) string {
	words := strings.Split(paramsBody, " ")

	cleanBody := []string{}

	for _, word := range words {
		if _, exists := profaneWords[strings.ToLower(word)]; exists {
			cleanBody = append(cleanBody, "****")
		} else {
			cleanBody = append(cleanBody, word)
		}
	}

	return strings.Join(cleanBody, " ")
}