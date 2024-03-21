package main

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func (apiCfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := strconv.Atoi(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error converting pathvalue to int")
	}
	tokenString := r.Header.Get("Authorization")
	trimmedToken := strings.TrimPrefix(tokenString, "Bearer ")
	claims := &jwt.RegisteredClaims{}
	_, err = jwt.ParseWithClaims(trimmedToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(apiCfg.JWTSecret), nil
	})
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "error converting pathvalue to int")
		return
	}
	err = apiCfg.DB.deleteChirp(chirpID, userID)
	if err != nil {
		respondWithError(w, http.StatusForbidden, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
