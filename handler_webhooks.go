package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/BoaPi/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhook(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	type response struct {
		User
	}

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "invalid auth header", err)
		return
	}

	if apiKey != cfg.polkaKey {
		responseWithError(w, http.StatusUnauthorized, "wrong api key", nil)
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	dat := request{}
	err = decoder.Decode(&dat)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode request", err)
		return
	}

	if dat.Event != "user.upgraded" {
		responseWithJSON(w, http.StatusNoContent, nil)
		return
	}

	_, err = cfg.db.UpgradeUserToChirpyRedById(r.Context(), dat.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			responseWithError(w, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		responseWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	responseWithJSON(w, http.StatusNoContent, nil)
}
