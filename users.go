package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (apiCfg *apiConfig) handlerUserCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldnt retrieve params")
		return
	}

	validated, err := validateUser(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := apiCfg.DB.CreateUser(validated)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:    user.ID,
		Email: user.Email,
	})

}

func validateUser(email string) (string, error) {
	const maxEmailLength = 140
	if len(email) > maxEmailLength {
		return "", errors.New("email is too long")
	}

	return email, nil
}
