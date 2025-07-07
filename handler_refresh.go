package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/rdcassin/server-go/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Token string `json:"token"`
	}

	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error fetching token in handlerRefresh: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), tokenString)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Error in handlerRefresh... token not in database: %s", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		log.Printf("Error fetching token in handlerRefresh: %s", err)
		respondWithInternalServerError(w)
		return
	}

	newJWTToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, auth.JWTTokenExpiration)
	if err != nil {
		log.Printf("Error creating new token in handlerRefresh: %s", err)
		respondWithInternalServerError(w)
		return
	}

	payload := returnVals{
		Token: newJWTToken,
	}

	respondWithJSON(w, http.StatusOK, payload)
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error fetching token in handlerRevoke: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = cfg.db.RevokeToken(r.Context(), tokenString)
	if err != nil {
		log.Printf("Error revoking token in handlerRevoke: %s", err)
		respondWithInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
