package main

import (
	"encoding/json"
	"net/http"

	"github.com/BoaPi/chirpy/internal/auth"
	"github.com/BoaPi/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	accesToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "missing auth header", err)
		return
	}

	userID, err := auth.ValidateJWT(accesToken, cfg.jwtSecret)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "couldn't validate JWT", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	dat := request{}
	err = decoder.Decode(&dat)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "couldn't decode request", err)
		return
	}

	hashedPassword, err := auth.HashPassword(dat.Password)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "couldn't hash new password", err)
		return
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             userID,
		Email:          dat.Email,
		HashedPassword: hashedPassword,
	})

	responseWithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		Email:       user.Email,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		IsChirpyRed: user.IsChirpyRed,
	})
}
