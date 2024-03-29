package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type User struct {
	Password      string `json:"password"`
	Email         string `json:"email"`
	ID            int    `json:"id"`
	Is_chirpy_red bool   `json:"is_chirpy_red"`
}

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
	if claims.Issuer != "chirpy-access" {
		respondWithError(w, 401, "wrong token type")
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

	err = apiCfg.DB.UpdateUserEmailPassword(userID, params.Email, params.Password)
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
		ID            int    `json:"id"`
		Email         string `json:"email"`
		Is_chirpy_red bool   `json:"is_chirpy_red"`
	}

	toReturn := ToReturn{
		ID:            userToReturn.ID,
		Email:         userToReturn.Email,
		Is_chirpy_red: userToReturn.Is_chirpy_red,
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
	type Response struct {
		ID            int    `json:"id"`
		Email         string `json:"email"`
		Is_chirpy_red bool   `json:"is_chirpy_red"`
	}
	response := Response{
		ID:            user.ID,
		Email:         user.Email,
		Is_chirpy_red: user.Is_chirpy_red,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func validateUser(email string) (string, error) {
	const maxEmailLength = 140
	if len(email) > maxEmailLength {
		return "", errors.New("email is too long")
	}

	return email, nil
}
