package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	type returnVals struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params in handlerCreateUser: %s", err)
		msg := "Something went wrong"
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	newUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating new user: %s", err)
		msg := "Error creating new user"
		respondWithError(w, http.StatusInternalServerError, msg)
		return
	}

	payload := returnVals{
		User: User{
			ID:        newUser.ID,
			CreatedAt: newUser.CreatedAt,
			UpdatedAt: newUser.UpdatedAt,
			Email:     newUser.Email,
		},
	}

	respondWithJSON(w, http.StatusCreated, payload)
}
