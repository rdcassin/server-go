package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		log.Fatal("Error resetting... Can only reset in development")
		msg := ""
		respondWithError(w, http.StatusForbidden, msg)
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		log.Fatalf("Error deleting all users: %s", err)
	}
}