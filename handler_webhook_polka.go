package main

import (
	"encoding/json"
	"net/http"

	"github.com/Kickjaw/HTTPServerProject/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) upgradeUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		}
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Error getting api key from headers", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithJSON(w, 401, nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldnt decode incoming body", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 204, nil)
		return
	}

	userUUID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error parsing user uuid", err)
		return
	}

	err = cfg.db.UpgradeUser(r.Context(), userUUID)
	if err != nil {
		respondWithError(w, 404, "error upgrading user in db", err)
		return
	}

	respondWithJSON(w, 204, nil)
}
