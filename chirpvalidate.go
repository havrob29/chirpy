package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type returnError struct {
	Error string `json:"error"`
}

type parameters struct {
	Body string `json:"body"`
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
	type returnValid struct {
		Body bool `json:"valid"`
	}
	respBody := returnValid{
		Body: true,
	}
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
