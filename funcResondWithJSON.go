package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	data, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, 500, fmt.Sprintf("Error marshalling JSON: %s", err))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)
}
