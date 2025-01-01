package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))
	apiCFG := apiConfig{
		fileserverHits: atomic.Int32{},
	}

	mux.Handle("/app/", apiCFG.middlewareMetricInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("GET /api/healthz", health)
	mux.HandleFunc("GET /admin/metrics", apiCFG.severMetrics)
	mux.HandleFunc("POST /admin/reset", apiCFG.resetServerHits)
	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

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
