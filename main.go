package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"context"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	fmt.Println("lfg!!")
	godotenv.Load()

	portString := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	log.Println(dbUrl)

	if portString == "" {
		log.Fatal("Port is not present")
	}

	if dbUrl == "" {
		log.Fatal("No DB URL found!!")
	}
	// db mongo db
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUrl))
	if client == nil {
		log.Fatal("MongoDB client is not initialized")
	}

	if err != nil {
		log.Fatal("Cannot connect to DB")
	}

	// ping the mongodb instance
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Printf("Mongo DB connection found at %v", dbUrl)

	router := chi.NewRouter()

	// cors configuration
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https//*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	// app client
	app := &App{Client: client}
	// v1Router.HandleFunc("/healthz", handlerReadiness)
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Get("/user/{userID}", app.handlerGetUserDetails)
	v1Router.Post("/create-user", app.handlerCreateUser)

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Server starting on port %v", portString)
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Port", portString)
}
