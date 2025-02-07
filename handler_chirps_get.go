package main

import (
	"net/http"
	"slices"
	"strings"

	"github.com/BoaPi/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	authorID := uuid.Nil
	authorIDString := r.URL.Query().Get("author_id")
	if authorIDString != "" {
		parsedAuthorID, err := uuid.Parse(authorIDString)
		if err != nil {
			responseWithError(w, http.StatusInternalServerError, "couldn't parse author id", err)
			return
		}
		authorID = parsedAuthorID
	}

	sortDirection := "asc"
	sortDirectionString := r.URL.Query().Get("sort")
	if sortDirectionString == "desc" {
		sortDirection = sortDirectionString
	}

	dbChirps, err := cfg.db.GetChirps(r.Context(), authorID)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "couldn't fetch chirps", err)
		return
	}

	if sortDirection == "desc" {
		slices.SortFunc(dbChirps, func(a database.Chirp, b database.Chirp) int {
			if a.CreatedAt.Before(b.CreatedAt) {
				return 1
			} else if a.CreatedAt.After(b.CreatedAt) {
				return -1
			}
			return 0
		})
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			UserID:    dbChirp.UserID,
			Body:      dbChirp.Body,
		})
	}

	responseWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerGetChirpById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("chirpID")
	if id == "" {
		responseWithError(w, http.StatusBadRequest, "no id given", nil)
		return
	}

	chirpID, err := uuid.Parse(id)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "invalid id", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		responseWithError(w, http.StatusNotFound, "couldn't get chirp", err)
		return
	}

	responseWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}
