package main

import (
	"fmt"
	"net/http"

	"github.com/AnkitNayan83/go-rss-agg/internal/auth"
	"github.com/AnkitNayan83/go-rss-agg/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, apiKey_err := auth.GetApiKey(r.Header)

		if apiKey_err != nil {
			respondWithError(w, 401, fmt.Sprintf("Api Key Error: %v", apiKey_err))
			return
		}

		user, user_err := cfg.DB.GetUserByApiKey(r.Context(), apiKey)

		if user_err != nil {
			respondWithError(w, 404, fmt.Sprintf("User not found: %v", user_err))
			return
		}

		handler(w, r, user)
	}
}
