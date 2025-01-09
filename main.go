package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Kickjaw/HTTPServerProject/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
	polkaKey       string
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	JWTSecret := os.Getenv("SECRET")
	if JWTSecret == "" {
		log.Fatal("JWT secrect must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connect to the database: %s", err)
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA key must be set")
	}
	dbQueries := database.New(dbConn)

	apiCFG := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		secret:         JWTSecret,
		polkaKey:       polkaKey,
	}

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))

	mux.Handle("/app/", apiCFG.middlewareMetricInc(http.StripPrefix("/app", fileServer)))

	mux.HandleFunc("GET /admin/metrics", apiCFG.severMetrics)
	mux.HandleFunc("POST /admin/reset", apiCFG.resetServerHits)

	mux.HandleFunc("POST /api/users", apiCFG.addUser)
	mux.HandleFunc("PUT /api/users", apiCFG.changeUserDetails)
	mux.HandleFunc("GET /api/healthz", health)

	mux.HandleFunc("POST /api/login", apiCFG.login)

	mux.HandleFunc("POST /api/refresh", apiCFG.refresh)
	mux.HandleFunc("POST /api/revoke", apiCFG.revoke)

	mux.HandleFunc("POST /api/chirps", apiCFG.Chirp)
	mux.HandleFunc("GET /api/chirps", apiCFG.GetChirpHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCFG.getChirpByID)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCFG.deleteChirpByID)

	mux.HandleFunc("POST /api/polka/webhooks", apiCFG.upgradeUser)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
