package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AnkitNayan83/go-rss-agg/internal/database"
	"github.com/google/uuid"
)

func (apiCfg apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if params.Name == "" {
		respondWithError(w, 400, "Name of the user cannot be empty")
	}

	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to create feed: %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedToFeed(feed))
}

func (apiCfg apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {

	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error in getting feeds: %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedsToFeeds(feeds))
}
