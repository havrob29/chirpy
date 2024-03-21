package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserWithoutPassword struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (apiCfg *apiConfig) putApiUser(w http.ResponseWriter, r *http.Request) {
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

	type RequestParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	params := RequestParams{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = apiCfg.DB.UpdateUser(userID, params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt update database")
		return
	}

	userArray, err := apiCfg.DB.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	userToReturn := User{}
	for _, user := range userArray {
		if user.ID == userID {
			userToReturn = user
		}
	}
	type ToReturn struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	toReturn := ToReturn{
		ID:    userToReturn.ID,
		Email: userToReturn.Email,
	}

	respondWithJSON(w, 200, toReturn)
}

func (apiCfg *apiConfig) postApiUsers(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	user := User{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt retrieve params")
		return
	}

	user.Email, err = validateUser(user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err = apiCfg.DB.CreateUser(user.Email, user.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, UserWithoutPassword{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (apiCfg *apiConfig) postApiLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type RequestParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	requestParams := RequestParams{}
	err := decoder.Decode(&requestParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt retrieve params")
		return
	}

	compareUser := User{}
	users, err := apiCfg.DB.GetUsers()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	//get user from database
	for _, user := range users {
		if user.Email == requestParams.Email {
			compareUser = user
		}
	}
	if compareUser.Email == "" {
		respondWithError(w, http.StatusUnauthorized, "email not in database")
		return
	}

	requestPassword := requestParams.Password
	savedPassword := compareUser.Password

	err = comparePassword(requestPassword, savedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	//if expiry time is not specified, or set to longer than 24 hours; set it to 24 hours

	accessJWTseconds := 3600
	refreshJWTseconds := 86400 * 60

	//create a access *JWT.Token for user ID
	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(accessJWTseconds))),
		Subject:   strconv.Itoa(compareUser.ID),
	})

	//create a refresh *JWT.Token for user ID
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-refresh",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Second * time.Duration(refreshJWTseconds))),
		Subject:   strconv.Itoa(compareUser.ID),
	})

	refreshToken, err := refreshJWT.SignedString([]byte(apiCfg.JWTSecret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	accessToken, err := accessJWT.SignedString([]byte(apiCfg.JWTSecret))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type Response struct {
		ID           int    `json:"id"`
		Email        string `json:"email"`
		AccessToken  string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	response := Response{
		ID:           compareUser.ID,
		Email:        compareUser.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	respondWithJSON(w, http.StatusOK, response)

}

func validateUser(email string) (string, error) {
	const maxEmailLength = 140
	if len(email) > maxEmailLength {
		return "", errors.New("email is too long")
	}

	return email, nil
}
