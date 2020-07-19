package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Config the configuration JSON structure
type Config struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

var config Config

func main() {
	// Configuration setup
	file, _ := ioutil.ReadFile("./config.json")
	json.Unmarshal(file, &config)

	// API Setup
	router := mux.NewRouter()

	// Routes
	// -------------------------------
	router.HandleFunc("/", Home)
	router.HandleFunc("/ping", Ping)
	// -------------------------------

	// Run server
	fmt.Printf("Server running on %s%s\n", config.Host, config.Port)
	log.Fatal(http.ListenAndServe(config.Port, router))

}
