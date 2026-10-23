package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) HandlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits.Load())))
}

func (cfg *apiConfig) HandlerReset(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Swap(0)
}

func (cfg *apiConfig) HandlerValidate(w http.ResponseWriter, req *http.Request) {
	// Expect JSON of {"body": "text here"}
	// If body is >140 characters then respond with JSON {"error": "Chirp is too long"} and code 400
	// If a different error JSON {"error": "Something went wrong"} and code 400
	// Otherwise JSON {"valid": true} and code 200

	type parameters struct {
		Body string `json: "body"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}
	if len(params.Body) > 140 {
		RespondWithError(w, 400, fmt.Sprintf("Chirp is too long"))
		return
	}

	type returnValid struct {
		Valid bool `json:"valid"`
	}
	resp := returnValid{
		Valid: true,
	}

	RespondWithJSON(w, 200, resp)
}
