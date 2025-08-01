package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithIncorrectEmailOrPassword(w http.ResponseWriter) {
	msg := "incorrect email or password"
	respondWithError(w, http.StatusUnauthorized, msg)
}

func respondWithInternalServerError(w http.ResponseWriter) {
	msg := "Something went wrong"
	respondWithError(w, http.StatusInternalServerError, msg)
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

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling data for response: %s", err)
		return
	}

	w.Write(data)
}