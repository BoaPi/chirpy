package main

import (
	"net/http"

	"github.com/BoaPi/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteChirpById(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, http.StatusUnauthorized, "couldn't find JWT", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "couldn't validate JWT", err)
		return
	}

	chirpIDString := r.PathValue("chirpID")
	if chirpIDString == "" {
		responseWithError(w, http.StatusBadRequest, "no id given", nil)
		return
	}
	chirpID, err := uuid.Parse(chirpIDString)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "invalid chirp id", err)
		return
	}

	dbChirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "couldn't find chirp", err)
		return
	}

	if dbChirp.UserID != userID {
		responseWithError(w, http.StatusForbidden, "not allowed to delete chirp", nil)
		return
	}

	err = cfg.db.DeleteChirpById(r.Context(), chirpID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "couldn't delete chirp", err)
		return
	}

	responseWithJSON(w, http.StatusNoContent, "chirp deleted")
}
