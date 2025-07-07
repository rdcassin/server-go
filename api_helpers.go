package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func decodeJSONBody(w http.ResponseWriter, r *http.Request, out interface{}) bool {
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(out); err != nil {
        log.Printf("Error decoding JSON body: %s", err)
        respondWithInternalServerError(w)
        return false
    }
    return true
}