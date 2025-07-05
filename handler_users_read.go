package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/rdcassin/server-go/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	payload := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	respondWithJSON(w, http.StatusOK, payload)
}
