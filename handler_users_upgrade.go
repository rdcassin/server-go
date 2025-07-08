package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/rdcassin/server-go/internal/auth"
)

func (cfg *apiConfig) handlerUsersUpgrade(w http.ResponseWriter, r *http.Request) {
	type UserIDFromPolka struct {
		UserID uuid.UUID `json:"user_id"`
	}

	type parameters struct {
		Event string          `json:"event"`
		Data  UserIDFromPolka `json:"data"`
	}

	receivedAPIKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error retrieving APIKey: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if receivedAPIKey != cfg.polkaKey {
		log.Print("API Key mismatch")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	params := parameters{}
	if !decodeJSONBody(w, r, &params) {
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpgradeUserByID(r.Context(), params.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
