package main

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Kickjaw/HTTPServerProject/internal/auth"
	"github.com/Kickjaw/HTTPServerProject/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

func (cfg *apiConfig) Chirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
		return
	}
	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}
	type response struct {
		Chirp
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "couldnt decode chirp", err)
		return
	}
	const maxChirpLength = 140

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is to long", nil)
		return
	}

	cleanedBody := removeProfanity(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "UUID is incorrect", err)
		return
	}

	chirpRaw, err := cfg.db.WriteChirpToDB(r.Context(), database.WriteChirpToDBParams{
		Body:   cleanedBody,
		UserID: userID,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error writing chirp to db", err)
		return
	}

	respondWithJSON(w, 201, response{
		Chirp: Chirp{
			ID:        chirpRaw.ID,
			CreatedAt: chirpRaw.CreatedAt,
			UpdatedAt: chirpRaw.UpdatedAt,
			Body:      chirpRaw.Body,
			UserID:    chirpRaw.UserID.String(),
		},
	})
}

func removeProfanity(Body string) string {
	replacer := []string{"Kerfuffle", "kerfuffle", "Sharbert", "sharbert", "Fornax", "fornax"}

	for _, profane := range replacer {
		Body = strings.ReplaceAll(Body, profane, "****")
	}

	return Body
}

func (cfg *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	sortOrder := r.URL.Query().Get("sort")
	chirps, err := cfg.db.RetrieveChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error retrieving chirps from db", err)
	}

	chirpsResponse := []Chirp{}

	for _, chirp := range chirps {
		chirpFormatted := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		}
		chirpsResponse = append(chirpsResponse, chirpFormatted)
	}

	if sortOrder == "desc" {
		sort.Slice(chirpsResponse, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, 200, chirpsResponse)

}

func (cfg *apiConfig) getChirpByID(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Chirp
	}
	chripId := r.PathValue("chirpID")
	if chripId == "" {
		w.WriteHeader(404)
		return
	}
	chirpUUID, err := uuid.Parse(chripId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error converting request id to UUID", err)
		return
	}
	chirpRaw, err := cfg.db.RetrieveByIDChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, 404, "chirp not found", err)
		return
	}

	respondWithJSON(w, 200, response{Chirp: Chirp{
		ID:        chirpRaw.ID,
		CreatedAt: chirpRaw.CreatedAt,
		UpdatedAt: chirpRaw.UpdatedAt,
		Body:      chirpRaw.Body,
		UserID:    chirpRaw.UserID.String(),
	},
	})
}

func (cfg *apiConfig) deleteChirpByID(w http.ResponseWriter, r *http.Request) {

	chripId := r.PathValue("chirpID")
	if chripId == "" {
		w.WriteHeader(404)
		return
	}

	chirpUUID, err := uuid.Parse(chripId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error converting request id to UUID", err)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "error getting bearer token", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
		return
	}

	chirpRaw, err := cfg.db.RetrieveByIDChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, 404, "chirp not found", err)
		return
	}

	if chirpRaw.UserID != userID {
		respondWithError(w, 403, "user not creator of chirp", err)
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpUUID)
	if err != nil {
		respondWithError(w, 404, "error deleting chirp", err)
		return
	}

	respondWithJSON(w, 204, nil)

}

func (cfg *apiConfig) GetChirpHandler(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	if authorID == "" {
		cfg.getAllChirps(w, r)
		return
	}
	sortOrder := r.URL.Query().Get("sort")
	authorUUID, err := uuid.Parse(authorID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error converting author id to UUID", err)
		return
	}

	chirps, err := cfg.db.RetrieveByAuthor(r.Context(), authorUUID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error finding chirps in db", err)
		return
	}
	chirpsResponse := []Chirp{}
	for _, chirp := range chirps {
		chirpsResponse = append(chirpsResponse, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID.String(),
		})
	}
	if sortOrder == "desc" {
		sort.Slice(chirpsResponse, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, 200, chirpsResponse)
}
