package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elue-dev/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

type feedParameters struct {
	Name string `json:"name"`
	Url string `json:"url"`
}


func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	decoder := json.NewDecoder(r.Body)

	params := feedParameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing json %v", err))
		return
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
		Name: params.Name,
		Url: params.Url,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not create feed %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedToFeedModel(feed))
}

func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Could not get feeds %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedsToFeedsModel(feeds))
}

