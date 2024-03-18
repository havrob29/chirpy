package main

import "strings"

func cleanBadWords(input string) string {
	badwords := map[string]string{
		"kerfuffle": "bad",
		"sharbert":  "bad",
		"fornax":    "bad",
	}
	words := strings.Split(input, " ")
	for i, word := range words {

		if badwords[word] == "bad" || badwords[strings.ToLower(word)] == "bad" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")
}
