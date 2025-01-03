package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Kickjaw/HTTPServerProject/internal/auth"
	"github.com/Kickjaw/HTTPServerProject/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) resetServerHits(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	err := cfg.db.DeleteUser(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error deleting users from db", err)
	}
	cfg.fileserverHits.Store(int32(0))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))

}

func (cfg *apiConfig) addUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldnt decode email", err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldnt hash password", err)
	}
	type CreateUserParams struct {
		Email          string
		HashedPassword string
	}

	userParams := database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hashedPass,
	}

	user, err := cfg.db.CreateUser(r.Context(), userParams)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error writing user to db", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}

func (cfg *apiConfig) login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldnt decode email or password", err)
		return
	}
	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, 401, "", err)
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "", err)
	}
	respondWithJSON(w, 200, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
