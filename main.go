package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/AnkitNayan83/go-rss-agg/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	env_err := godotenv.Load()

	if env_err != nil {
		log.Fatal("Error in loading .env file")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT not found in the environment")
	}

	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("Database URL not found in the environment")
	}

	conn, conn_err := sql.Open("postgres", dbUrl)

	if conn_err != nil {
		log.Fatal("Databese connection failed: ", conn_err)
	}

	queries := database.New(conn) // to convert conn type to database.queries type

	apiCfg := apiConfig{
		DB: queries,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))
	v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	v1Router.Post("/feed-follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollows))

	router.Mount("/api/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v\n", portString)

	server_err := srv.ListenAndServe()
	if server_err != nil {
		log.Fatal(server_err)
	}

}
