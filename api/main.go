package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alexander-cho/manager/api/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("Hello there.")

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}

	fmt.Println("Listening at:", portString)

	// import database connection
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB URL is not found is not found in the environment")
	}

	// go std lib sql package
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}

	// create database connection
	queries := database.New(conn)

	// create new api config, pass into handlers to give access to database
	apiCfg := apiConfig{
		DB: queries,
	}

	// define router
	router := chi.NewRouter()

	// cors config for router, let user access from browser
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           600,
	}))

	v1Router := chi.NewRouter()
	// hook up functions to "/ready" path with request methods
	v1Router.Get("/ready", handlerReadiness)
	v1Router.Get("/error", handlerError)
	v1Router.Post("/users", apiCfg.handlerCreateUser)
	v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUser))

	v1Router.Post("/entries", apiCfg.middlewareAuth(apiCfg.handlerCreateEntry))
	v1Router.Get("/entries", apiCfg.handlerGetEntries)

	v1Router.Post("/following-entries", apiCfg.middlewareAuth(apiCfg.handlerCreateFollowingEntry))
	v1Router.Get("/following-entries", apiCfg.middlewareAuth(apiCfg.handlerGetFollowingEntries))
	v1Router.Delete("/following-entries/{followingEntryID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFollowingEntry))

	// full path: /v1/ready
	router.Mount("/v1", v1Router)

	// define server
	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server starting on port %v", portString)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
