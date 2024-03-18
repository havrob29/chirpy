package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type returnError struct {
	Error string `json:"error"`
}

type parameters struct {
	Body string `json:"body"`
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON. %s", err)
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respBody := returnError{
		Error: msg,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON. %s", err)
		w.WriteHeader(500)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func validateChirp(w http.ResponseWriter, r *http.Request) {

	//decode json request body

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	//if chirp > 140 characters, invalidate chirp and
	//return json with error message:
	if len(params.Body) >= 140 {

		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
		//else if chirp valid:
	}

	params.Body = cleanBadWords(params.Body)

	type cleanedStruct struct {
		Body string `json:"cleaned_body"`
	}
	respBody := cleanedStruct{
		Body: params.Body,
	}
	respondWithJSON(w, 200, respBody)
}

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
