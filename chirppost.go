package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Chirp struct {
	Body   string `json:"body"`
	ID     int    `json:"id"`
	Author int    `json:"author_id"`
}

func (apiCfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	trimmedToken := strings.TrimPrefix(tokenString, "Bearer ")
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(trimmedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(apiCfg.JWTSecret), nil
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt retrieve params")
		return
	}

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	chirp, err := apiCfg.DB.CreateChirp(cleaned, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)

}

func validateChirp(body string) (string, error) {
	const maxChirpLength = 140
	if len(body) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	cleaned := cleanBadWords(body)
	return cleaned, nil
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
