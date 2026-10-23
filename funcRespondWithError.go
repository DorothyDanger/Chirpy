package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	type returnError struct {
		Error string `json:"error"`
	}
	respBody := returnError{
		Error: msg,
	}
	data, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
	}
	w.WriteHeader(code)
	w.Write(data)
}
