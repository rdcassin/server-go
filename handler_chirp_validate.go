package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func handlerChirpValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		msg := "Something went wrong"
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	chars := []rune(params.Body)
	if len(chars) > 140 {
		msg := "Chirp is too long"
		respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	type successResponse struct {
		Valid bool `json:"valid"`
	}

	payload := successResponse {
		Valid: true,
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	payload := errorResponse {
		Error: msg,
	}

	respondWithJSON(w, code, payload)
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