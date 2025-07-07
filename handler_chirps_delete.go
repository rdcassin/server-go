package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	userID, err := cfg.validateUser(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	
	chirp, err := cfg.getChirp(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if userID != chirp.UserID {
		log.Print("Error deleting chirp... user not chirp author")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirp.ID)
	if err != nil {
		log.Printf("Error deleting chirp ID %s: %s", chirp.ID, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}