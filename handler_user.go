package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elue-dev/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

type parameters struct {
	Name string `json:"name"`
}

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json %v", err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		Name: params.Name,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create user %v", err))
		return
	}

	respondWithJSON(w, 201, databaseUserToUserModel(user))
}

func (apiCfg *apiConfig) handlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUserToUserModel(user))
}