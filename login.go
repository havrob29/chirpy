package main

import (
	"encoding/json"
	"net/http"
)

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

	accessToken, err := apiCfg.makeAccessToken(compareUser.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	refreshToken, err := apiCfg.makeRefreshToken(compareUser.ID)
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
