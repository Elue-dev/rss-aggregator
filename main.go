package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/elue-dev/rss-aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type errResponse struct {
	Error string `json:"error"`
}

type apiConfig struct {
	DB *database.Queries
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Println("Responding with 5XX error:", msg)
	}

	respondWithJSON(w, code, errResponse{
		Error: msg,
	})
}

func main() {
	_, err := urlToFeed("https://blog.boot.dev/index.xml")
	if err != nil {
		log.Fatal(err)
	}
	
	godotenv.Load(".env")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("port could not be found in env")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL could not be found in env")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database", err)
	}

	db := database.New(conn)
	apiCfg := apiConfig{
		DB:  db,
	}

	go startScrapping(db, 10, time.Minute)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, 
	  }))

	  v1Router := chi.NewRouter()

	  v1Router.Get("/healthz", handlerReadiness)
	  v1Router.Get("/err", hanndlerErr)
	  v1Router.Post("/users", apiCfg.handlerCreateUser)
	  v1Router.Get("/users", apiCfg.middlewareAuth(apiCfg.handlerGetUserByAPIKey))

	  v1Router.Get("/posts", apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

	  v1Router.Post("/feeds", apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
	  v1Router.Get("/feeds", apiCfg.handlerGetFeeds)
	  v1Router.Post("/feeds_follows", apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
	  v1Router.Get("/feeds_follows", apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
	  v1Router.Delete("/feeds_follows/{feedFollowID}", apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollow))
	  

	  router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr: ":" + port,
	}

	log.Printf("Server starting on port %v", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}