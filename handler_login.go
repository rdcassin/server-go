package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/rdcassin/server-go/internal/auth"
	"github.com/rdcassin/server-go/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) handlerUsersLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type returnVals struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := parameters{}
	if !decodeJSONBody(w, r, &params) {
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Error due to incorrect user/password: %s", err)
			respondWithIncorrectEmailOrPassword(w)
			return
		}
		log.Printf("Error fetching user by email in handlerUsersLogin: %s", err)
		respondWithInternalServerError(w)
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword.String)
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			log.Printf("Error due to incorrect user/password: %s", err)
			respondWithIncorrectEmailOrPassword(w)
			return
		}
		log.Printf("Error checking password in handlerUsersLogin: %s", err)
		respondWithInternalServerError(w)
		return
	}

	// Create new token
	newJWTToken, err := auth.MakeJWT(user.ID, cfg.tokenSecret, auth.JWTTokenExpiration)
	if err != nil {
		log.Printf("Error creating new token in handlerUsersLogin: %s", err)
		respondWithInternalServerError(w)
		return
	}

	// Create new refresh token
	newRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("Error creating new refresh token in handlerUsersLogin: %s", err)
		respondWithInternalServerError(w)
		return
	}

	// Load new refresh token in refresh_tokens DB
	refreshTokenParams := database.CreateRefreshTokenParams{
		Token:     newRefreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(auth.RefreshTokenExpiration),
	}

	refreshToken, err := cfg.db.CreateRefreshToken(r.Context(), refreshTokenParams)
	if err != nil {
		log.Printf("Error creating new refresh token in database in handlerUsersLogin: %s", err)
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
		Token:        newJWTToken,
		RefreshToken: refreshToken.Token,
	}

	respondWithJSON(w, http.StatusOK, payload)
}
