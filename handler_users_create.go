package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type createUserRequest struct {
	Email string `json:"email"`
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User
	}

	decoder := json.NewDecoder(r.Body)
	req := createUserRequest{}
	err := decoder.Decode(&req)
	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't decode request", err)
		return
	}
	defer r.Body.Close()

	user, err := cfg.db.CreateUser(r.Context(), req.Email)

	if err != nil {
		responseWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	responseWithJSON(w, http.StatusCreated, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
