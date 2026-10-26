package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/DorothyDanger/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
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
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		RespondWithError(w, 403, "Forbidden")
	}
	err := cfg.db.ResetUsers(req.Context())
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Error resetting users: %s", err))
	}
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

	// For lesson 4 in JSON: Replace "profane" words with ****
	// The words are:
	// kerfuffle
	// sharbert
	// fornax
	// instead of the valid boolean, return the cleaned string in a JSON response.

	cleanedBody := ProfaneReplacement(params.Body)
	type returnBody struct {
		CleanedBody string `json:"cleaned_body"`
	}
	resp := returnBody{
		CleanedBody: cleanedBody,
	}

	RespondWithJSON(w, 200, resp)
}

func (cfg *apiConfig) HandlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Email string `json: "email"`
	}

	decoder := json.NewDecoder(req.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}
	var finalUser User
	user, err := cfg.db.CreateUser(req.Context(), params.Email)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Error creating user: %s", err))
		return
	}
	finalUser = User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	RespondWithJSON(w, 201, finalUser)
}

func (cfg *apiConfig) HandlerDeleteAllUsers(w http.ResponseWriter, req *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		RespondWithError(w, 403, "Forbidden")
	}
	err := cfg.db.ResetUsers(req.Context())
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Error resetting users: %s", err))
	}
}
