package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersUpgrade(w http.ResponseWriter, r *http.Request) {
	type UserIDFromPolka struct {
		UserID uuid.UUID `json:"user_id"`
	}

	type parameters struct {
		Event string          `json:"event"`
		Data  UserIDFromPolka `json:"data"`
	}

	params := parameters{}
	if !decodeJSONBody(w, r, &params) {
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err := cfg.db.UpgradeUserByID(r.Context(), params.Data.UserID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	return
}
