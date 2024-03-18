package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/havrob29/chirpy/internal/database"
)

type returnError struct {
	Error string `json:"error"`
}

func getChirp(w http.ResponseWriter, r *http.Request) {

}

func (apiCfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	chirp := database.Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(chirp.Body) >= 140 {

		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	chirp.Body = cleanBadWords(chirp.Body)

	type structwithID struct {
		ID   int    `json:"id"`
		Body string `json:"body"`
	}
	respBody := structwithID{
		ID:   apiCfg.chirpCount,
		Body: chirp.Body,
	}
	respondWithJSON(w, 201, respBody)

	apiCfg.chirpCount++
}
