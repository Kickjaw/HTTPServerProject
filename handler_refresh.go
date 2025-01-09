package main

import (
	"net/http"
	"time"

	"github.com/Kickjaw/HTTPServerProject/internal/auth"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "error getting bearer token", err)
		return
	}
	User, err := cfg.db.FindRefreshToken(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "error finding refresh token", err)
		return
	}
	JWT, err := auth.MakeJWT(User.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 401, "error making jwt", err)
		return
	}

	respondWithJSON(w, 200, response{
		Token: JWT,
	})

}

func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "error getting token from header", err)
		return
	}

	err = cfg.db.RevokeRefreshToke(r.Context(), token)
	if err != nil {
		respondWithError(w, 401, "couldn't revoke session", err)
		return
	}
	respondWithJSON(w, 204, nil)
}
