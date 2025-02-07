package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/BoaPi/chirpy/internal/auth"
	"github.com/BoaPi/chirpy/internal/database"
)

type loginUserRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
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

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create access JWT", err)
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create refresh token entry", err)
	}

	responseWithJSON(w, http.StatusOK, response{
		User: User{
			ID:          user.ID,
			CreatedAt:   user.CreatedAt,
			UpdatedAt:   user.UpdatedAt,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
