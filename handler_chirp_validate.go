package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIResponse interface {
	isAPIResponse()
}

type errorResponse struct {
	Error string `json:"error"`
}

func (e errorResponse) isAPIResponse() {}

type successResponse struct {
	Valid bool `json:"valid"`
}

func (s successResponse) isAPIResponse() {}

func handlerChirpValidate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		bodyData := "Something went wrong"
		sendResponse(w, bodyData, 500)
		return
	}

	chars := []rune(params.Body)
	if len(chars) > 140 {
		bodyData := "Chirp is too long"
		sendResponse(w, bodyData, 400)
		return
	}

	bodyData := true
	sendResponse(w, bodyData, 200)
}

func sendResponse(w http.ResponseWriter, bodyData interface{}, statusCode int) {
	var respBody APIResponse

	if statusCode >= 400 {
		respBody = errorResponse{
			Error: bodyData.(string),
		}
	} else {
		respBody = successResponse{
			Valid: bodyData.(bool),
		}
	}

	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(data)
}
