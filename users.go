package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type UserWithoutPassword struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (apiCfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {

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

func (apiCfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	type RequestParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		Expires  int    `json:"expires_in_seconds"`
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

	respUser := UserWithoutPassword{
		ID:    compareUser.ID,
		Email: compareUser.Email,
	}

	respondWithJSON(w, http.StatusOK, respUser)

}

func validateUser(email string) (string, error) {
	const maxEmailLength = 140
	if len(email) > maxEmailLength {
		return "", errors.New("email is too long")
	}

	return email, nil
}
