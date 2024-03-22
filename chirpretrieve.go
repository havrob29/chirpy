package main

import (
	"net/http"
	"sort"
	"strconv"
)

func (apiCfg *apiConfig) getApiChirps(w http.ResponseWriter, r *http.Request) {
	s := r.URL.Query().Get("author_id")
	optionalAuthorOnly := false
	if s != "" {
		optionalAuthorOnly = true
	}

	dbChirps, err := apiCfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "cant retrieve chrips")
		return
	}
	chirps := []Chirp{}

	if !optionalAuthorOnly {
		for _, dbChirp := range dbChirps {
			chirps = append(chirps, Chirp{
				ID:     dbChirp.ID,
				Body:   dbChirp.Body,
				Author: dbChirp.Author,
			})
		}
	} else {
		authorID, err := strconv.Atoi(s)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "cant retrieve chrips")
			return
		}
		for _, dbChirp := range dbChirps {
			if dbChirp.Author == authorID {
				chirps = append(chirps, Chirp{
					ID:     dbChirp.ID,
					Body:   dbChirp.Body,
					Author: dbChirp.Author,
				})
			}
		}
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})
	respondWithJSON(w, http.StatusOK, chirps)
}

func (apiCfg *apiConfig) getApiChirpsByID(w http.ResponseWriter, r *http.Request) {

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
