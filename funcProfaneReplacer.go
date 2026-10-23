package main

import (
	"log"
	"strings"
)

// Helper function that takes a list of strings, breaks them up and replaces the profane words with ****
func ProfaneReplacement(body string) string {
	words := strings.Split(body, " ")
	filter := []string{"kerfuffle", "sharbert", "fornax"}
	for i, word := range words {
		for _, comp := range filter {
			log.Printf("\nWord: %s, Filter: %s\n", word, comp)
			if strings.Compare(strings.ToLower(word), strings.ToLower(comp)) == 0 {
				words[i] = "****"
				continue
			}
		}
	}
	cleaned := strings.Join(words, " ")
	return cleaned
}
