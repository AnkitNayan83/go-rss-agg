package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AnkitNayan83/go-rss-agg/internal/database"
	"github.com/google/uuid"
)

func (apiCfg apiConfig) handlerCreateFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	if params.FeedId == uuid.Nil {
		respondWithError(w, 400, "Name of the user cannot be empty")
	}

	feed_follows, err := apiCfg.DB.CreateFeedFollows(r.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to create feed follow: %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedFollowsToFeedFollows(feed_follows))
}

// func (apiCfg apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {

// 	feeds, err := apiCfg.DB.GetFeeds(r.Context())

// 	if err != nil {
// 		respondWithError(w, 500, fmt.Sprintf("Error in getting feeds: %v", err))
// 		return
// 	}

// 	respondWithJson(w, 201, databaseFeedsToFeeds(feeds))
// }
