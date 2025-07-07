package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/rdcassin/server-go/internal/auth"
)

func (cfg *apiConfig) validateUser(header http.Header) (uuid.UUID, error) {
	tokenString, err := auth.GetBearerToken(header)
	if err != nil {
		log.Printf("Error fetching bearer token in handlerAddChirp: %s", err)
		return uuid.Nil, err
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.tokenSecret)
	if err != nil {
		log.Printf("Error validation token... unauthorized in handlerAddChirp: %s", err)
		return uuid.Nil, err
	}

	return userID, nil
}
