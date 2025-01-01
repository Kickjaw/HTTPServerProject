package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}
	type returnVals struct {
		CleanBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := chirp{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	const maxChirpLength = 140

	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is to long", nil)
		return
	}

	cleanedBody := removeProfanity(params.Body)

	respondWithJSON(w, http.StatusOK, returnVals{
		CleanBody: cleanedBody,
	})

}

func removeProfanity(Body string) string {
	replacer := []string{"Kerfuffle", "kerfuffle", "Sharbert", "sharbert", "Fornax", "fornax"}

	for _, profane := range replacer {
		Body = strings.ReplaceAll(Body, profane, "****")
	}

	return Body
}
