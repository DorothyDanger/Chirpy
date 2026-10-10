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
	serveMux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	serveMux.HandleFunc("/healthz", handlerReadiness)
	log.Fatal(serverStruct.ListenAndServe())
}
