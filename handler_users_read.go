package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/rdcassin/server-go/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

// Max Expiration Time is set to 1 hour = 60 sec/min * 60 min/hour, or 3600 seconds
const maxExpiration = 3600

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type returnVals struct {
		User
		Token string `json:"token"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params in handlerUsersLogin: %s", err)
		respondWithInternalServerError(w)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error fetching user by email in handlerUsersLogin: %s", err)
		respondWithIncorrectEmailOrPassword(w)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword.String)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Printf("Error due to incorrect password: %s", err)
			respondWithIncorrectEmailOrPassword(w)
			return
		}
		log.Printf("Error checking password in hendlerUsersLogin: %s", err)
		respondWithInternalServerError(w)
		return
	}

	if err != nil {
		log.Fatalf("Error parsing User ID in handlerUsersLogin: %s", err)
		respondWithInternalServerError(w)
		return
	}

	expiration := time.Hour
	if params.ExpiresInSeconds < maxExpiration && params.ExpiresInSeconds > 0 {
		expiration = time.Duration(params.ExpiresInSeconds) * time.Second
	}

	// Create new token
	newToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, expiration)
	if err != nil {
		log.Fatalf("Error creating new token: %s", err)
		respondWithInternalServerError(w)
		return
	}

	payload := returnVals{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: newToken,
	}

	respondWithJSON(w, http.StatusOK, payload)
}
