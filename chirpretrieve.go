package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (apiCfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := apiCfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cant retrieve chrips")
		return
	}
	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}
	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (apiCfg *apiConfig) handlerSingleRetrieve(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "err converting string to int")
		return
	}
	dbChirps, err := apiCfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cant retrieve chrips")
		return
	}
	chirp := Chirp{}
	for _, dbChirp := range dbChirps {
		if dbChirp.ID == id {
			chirp = dbChirp
		}
	}
	if chirp.Body == "" && chirp.ID == 0 {
		respondWithError(w, 404, "chirp not found")
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
