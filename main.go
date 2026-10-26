package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/DorothyDanger/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	var apiCfg apiConfig
	apiCfg = apiConfig{
		db: dbQueries,
	}

	serveMux := http.NewServeMux()
	serverStruct := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	//serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)
	serveMux.HandleFunc("POST /api/validate_chirp", apiCfg.HandlerValidate)
	serveMux.HandleFunc("POST /api/users", apiCfg.HandlerCreateUser)
	//serveMux.HandleFunc("POST /admin/reset", apiCfg.HandlerDeleteAllUsers)
	log.Fatal(serverStruct.ListenAndServe())
}
