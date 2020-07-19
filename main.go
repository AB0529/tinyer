package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config the configuration JSON structure
type Config struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	MongoURI string `json:"mongouri"`
}

var config Config
var ctx context.Context
var cancel context.CancelFunc
var db *mongo.Collection

func main() {
	// Configuration setup
	file, _ := ioutil.ReadFile("./config.json")
	json.Unmarshal(file, &config)

	// MongoDB setup
	client, err := mongo.NewClient(options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		panic(err)
	}
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	// The URLS collections in database
	db = client.Database("moistdb").Collection("urls")
	fmt.Println("Mongo connection successful...")

	// API Setup
	router := mux.NewRouter()

	// Routes
	// -------------------------------

	router.HandleFunc("/", Home)
	router.HandleFunc("/urls/{id}", GetURL).Methods("GET")
	router.HandleFunc("/urls", CreateURL).Methods("POST")
	router.HandleFunc("/urls/{id}", DeleteURL).Methods("DELETE")
	router.HandleFunc("/{slug}", Redirect)

	// -------------------------------

	// Run server
	fmt.Printf("Server running on %s%s\n", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(config.Port, router))

}
