package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/iuliancarnaru/rss-aggregator/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	godotenv.Load(".env")

	portStr := os.Getenv("PORT")
	if portStr == "" {
		log.Fatal("PORT is not found in the environment")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	db := database.New(conn)
	apiConf := apiConfig{
		DB: db,
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	go startScraping(db, 10, time.Minute)

	v1Router := chi.NewRouter()

	// HEALTH
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/error", handlerError)

	// USERS
	v1Router.Post("/users", apiConf.handlerCreateUser)
	v1Router.Get("/users", apiConf.middlewareAuth(apiConf.handlerGetUser))

	// FEEDS
	v1Router.Post("/feeds", apiConf.middlewareAuth(apiConf.handlerCreateFeed))
	v1Router.Get("/feeds", apiConf.handlerGetFeeds)

	// FEED FOLLOWS
	v1Router.Post("/feed_follows", apiConf.middlewareAuth(apiConf.handlerCreateFeedFollows))
	v1Router.Get("/feed_follows", apiConf.middlewareAuth(apiConf.handlerGetFeedFollows))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiConf.middlewareAuth(apiConf.handlerDeleteFeedFollows))

	// POSTS
	v1Router.Get("/posts", apiConf.middlewareAuth(apiConf.handlerGetPostsForUser))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portStr,
	}

	log.Printf("Server starting on port %v", portStr)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
