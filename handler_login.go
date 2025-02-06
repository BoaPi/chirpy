package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/BoaPi/chirpy/internal/auth"
)

type loginUserRequest struct {
	Password         string `json:"password"`
	Email            string `json:"email"`
	ExpiresInSeconds int    `json:"expires_in_seconds"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	decoder := json.NewDecoder(r.Body)
	dat := loginUserRequest{}

	err := decoder.Decode(&dat)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode request", err)
		return
	}
	defer r.Body.Close()

	user, err := cfg.db.GetUserByEmail(r.Context(), dat.Email)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(dat.Password, user.HashedPassword)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "incorrect email or password", err)
		return
	}

	expiresTime := time.Hour
	if dat.ExpiresInSeconds > 0 && dat.ExpiresInSeconds < 3600 {
		expiresTime = time.Duration(dat.ExpiresInSeconds) * time.Second
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expiresTime)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}

	responseWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: accessToken,
	})
}
