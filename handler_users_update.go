package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/rdcassin/server-go/internal/auth"
	"github.com/rdcassin/server-go/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type returnVals struct {
		User
	}

	userID, err := cfg.validateUser(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := parameters{}
	decodeJSONBody(w, r, &params)

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		log.Printf("Error encrypting password in handlerUsersUpdate")
		respondWithInternalServerError(w)
		return
	}

	updatedUserParams := database.UpdateUserByIDParams{
		ID:    userID,
		Email: params.Email,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
	}

	updatedUser, err := cfg.db.UpdateUserByID(r.Context(), updatedUserParams)
	if err != nil {
		log.Printf("Error updating user password and email: %s", err)
		respondWithInternalServerError(w)
		return
	}

	payload := returnVals{
		User: User{
			ID:        updatedUser.ID,
			CreatedAt: updatedUser.CreatedAt,
			UpdatedAt: updatedUser.UpdatedAt,
			Email:     updatedUser.Email,
			IsChirpyRed: updatedUser.IsChirpyRed,
		},
	}

	respondWithJSON(w, http.StatusOK, payload)
}
