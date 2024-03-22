package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (apiCfg *apiConfig) postApiPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	authString := r.Header.Get("Authorization")
	ApiKey := strings.TrimPrefix(authString, "ApiKey ")

	if ApiKey != apiCfg.polka_key {
		respondWithError(w, http.StatusUnauthorized, "apikey mismatch")
		return
	}

	type RequestParams struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}
	params := RequestParams{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusOK)
		return
	}
	err = apiCfg.DB.UpgradeUser(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, nil)
}
