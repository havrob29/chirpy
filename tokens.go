package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// makes access token for userID with 1 hour duration
func (apiCfg *apiConfig) makeAccessToken(userID int) (token string, err error) {
	//set expiry time to 1hr
	accessJWTseconds := 3600
	// create a access *JWT.Token for user ID
	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(accessJWTseconds))),
		Subject:   strconv.Itoa(userID),
	})
	accessToken, err := accessJWT.SignedString([]byte(apiCfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return accessToken, nil
}

// create refresh token for users ID with 60 days duration
func (apiCfg *apiConfig) makeRefreshToken(userID int) (token string, err error) {
	refreshJWTseconds := 86400 * 60
	// create a refresh *JWT.Token for user ID
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(refreshJWTseconds))),
		Subject:   strconv.Itoa(userID),
	})
	refreshToken, err := refreshJWT.SignedString([]byte(apiCfg.JWTSecret))
	if err != nil {
		return "", err
	}
	return refreshToken, nil
}

func (apiCfg *apiConfig) postApiRefresh(w http.ResponseWriter, r *http.Request) {
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

	if claims.Issuer != "chirpy-refresh" {
		respondWithError(w, http.StatusUnauthorized, "not a refresh token")
		return
	}

	revokedTokens, err := apiCfg.DB.GetRevoked()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	for _, revokedToken := range revokedTokens {
		if revokedToken == trimmedToken {
			respondWithError(w, http.StatusUnauthorized, "token has been revoked")
			return
		}
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	newAccessToken, err := apiCfg.makeAccessToken(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type Response struct {
		Token string `json:"token"`
	}

	response := Response{
		Token: newAccessToken,
	}

	respondWithJSON(w, http.StatusOK, response)

}

func (apiCfg *apiConfig) postApiRevoke(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	trimmedToken := strings.TrimPrefix(tokenString, "Bearer ")
	err := apiCfg.DB.CreateRevoked(trimmedToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
