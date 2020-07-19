package server

import (
	"encoding/json"
	"net/http"
)

// Response the API response object
type Response struct {
	Status int
	State  string
	Result interface{}
}

// Home serves the home html file
func Home(w http.ResponseWriter, r *http.Request) {
	// TODO: Create better home page redirects for now
	http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusSeeOther)
}

// Ping test route to see if API is online
func Ping(w http.ResponseWriter, r *http.Request) {
	SendJSON(w, Response{Status: 200, State: "ok", Result: "Pong!"})
}

// SendJSON util func to handle sending JSON response
func SendJSON(w http.ResponseWriter, resp Response) (bool, error) {
	// Decode JSON
	json, err := json.Marshal(resp)

	if err != nil {
		return false, err
	}

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)

	return true, nil
}
