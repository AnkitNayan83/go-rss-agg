package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/AnkitNayan83/go-rss-agg/internal/database"
	"github.com/go-chi/chi"
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

	respondWithJson(w, 201, databaseFeedFollowToFeedFollow(feed_follows))
}

func (apiCfg apiConfig) handlerGetFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {

	feed_follows, err := apiCfg.DB.GetUserFeedFollows(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Error in getting feed follows: %v", err))
		return
	}

	respondWithJson(w, 201, databaseFeedFollowsToFeedFollows(feed_follows))
}

func (apiCfg apiConfig) handlerDeleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollowIDStr := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowIDStr)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Cannot parse feed follows id: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Cannot delete feed follows: %v", err))
		return
	}

	type deleteResponse struct {
		Message string `json:"message"`
	}

	respondWithJson(w, 200, deleteResponse{
		Message: "Feed unfollowed successfully",
	})

}
