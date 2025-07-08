package main

import (
	"database/sql"
	"github.com/rdcassin/server-go/internal/auth"
	"github.com/rdcassin/server-go/internal/database"
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type returnVals struct {
		User
	}

	params := parameters{}
	if !decodeJSONBody(w, r, &params) {
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
			Valid:  true,
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
			ID:          newUser.ID,
			CreatedAt:   newUser.CreatedAt,
			UpdatedAt:   newUser.UpdatedAt,
			Email:       newUser.Email,
			IsChirpyRed: newUser.IsChirpyRed,
		},
	}

	respondWithJSON(w, http.StatusCreated, payload)
}
