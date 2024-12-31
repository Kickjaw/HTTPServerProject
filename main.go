package main

import (
	"fmt"
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
	mux.HandleFunc("GET /healthz", health)
	mux.HandleFunc("GET /metrics", apiCFG.serverHits)
	mux.HandleFunc("POST /reset", apiCFG.resetServerHits)

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

func (cfg *apiConfig) middlewareMetricInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})

}

func (cfg *apiConfig) serverHits(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	hits := cfg.fileserverHits.Load()

	w.Write([]byte(fmt.Sprintf("Hits: %d", hits)))
}

func (cfg *apiConfig) resetServerHits(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(int32(0))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))

}
