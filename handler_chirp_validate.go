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

func handlerChirpValidate(w http.ResponseWriter, r *http.Request) {
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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	payload := errorResponse{
		Error: msg,
	}

	respondWithJSON(w, code, payload)
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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling data: %s", err)
		return
	}

	w.Write(data)
}
