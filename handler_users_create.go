package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/rdcassin/server-go/internal/auth"
	"github.com/rdcassin/server-go/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	type returnVals struct {
		User
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding params in handlerUsersCreate: %s", err)
		respondWithInternalServerError(w)
		return
	}

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error encrypting password in handlerUsersCreate")
		respondWithInternalServerError(w)
		return
	}

	newUserParams := database.CreateUserParams{
		Email: params.Email,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid: true,
		},
	}

	newUser, err := cfg.db.CreateUser(r.Context(), newUserParams)
	if err != nil {
		log.Printf("Error creating new user: %s", err)
		respondWithInternalServerError(w)
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
