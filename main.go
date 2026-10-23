package main

import (
	"log"
	"net/http"
)

func main() {

	serveMux := http.NewServeMux()
	serverStruct := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	var apiCfg apiConfig

	//serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	serveMux.HandleFunc("GET /api/healthz", handlerReadiness)
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.HandlerMetrics)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.HandlerReset)
	serveMux.HandleFunc("POST /api/validate_chirp", apiCfg.HandlerValidate)
	log.Fatal(serverStruct.ListenAndServe())
}
