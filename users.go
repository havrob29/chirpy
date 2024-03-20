package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

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

	type ResponseWithoutPassword struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	respondWithJSON(w, http.StatusCreated, ResponseWithoutPassword{
		ID:    user.ID,
		Email: user.Email,
	})

}

func (apiCfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	return
}

func validateUser(email string) (string, error) {
	const maxEmailLength = 140
	if len(email) > maxEmailLength {
		return "", errors.New("email is too long")
	}

	return email, nil
}
