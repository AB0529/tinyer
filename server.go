package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"unsafe"

	"github.com/gorilla/mux"
	s "github.com/gosimple/slug"
	"go.mongodb.org/mongo-driver/bson"
)

// Response the API response object
type Response struct {
	Status int         `json:"status"`
	State  string      `json:"state"`
	Result interface{} `json:"result"`
}

// Tinyer the url that was tinyified
type Tinyer struct {
	Slug      string `json:"slug"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	CreatedAt string `json:"created-at"`
}

// Home serves the home html file
func Home(w http.ResponseWriter, r *http.Request) {
	// TODO: Create better home page redirects for now
	http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusSeeOther)
}

// Redirect will redirect user if valid id is provided
func Redirect(w http.ResponseWriter, r *http.Request) {
	id := s.Make(mux.Vars(r)["slug"])
	// Validate identifier
	res, cancel, _ := FindSlug(bson.M{"slug": id})
	defer cancel()

	if res == nil {
		fmt.Fprint(w, "404 not found")
		return
	}

	if res["slug"] == id {
		// Redirect to url
		http.Redirect(w, r, fmt.Sprintf("%v", res["url"]), http.StatusSeeOther)
		return
	}
}

// GetURL will return info about the URL
func GetURL(w http.ResponseWriter, r *http.Request) {
	id := s.Make(mux.Vars(r)["id"])

	// Search Mongo with identifier
	res, cancel, _ := FindSlug(bson.M{"slug": id})
	defer cancel()

	if res == nil {
		SendJSON(w, Response{Status: http.StatusNotFound, State: "fail", Result: fmt.Sprintf("error: url with identifier '%s' could not be found", id)})
		return
	}

	if res["slug"] == id {
		// Display result
		SendJSON(w, Response{Status: http.StatusOK, State: "ok", Result: res})
		return
	}
}

// CreateURL creates a new url and uploads to database
func CreateURL(w http.ResponseWriter, r *http.Request) {
	var bod Tinyer
	var slug string

	// Decode incomming request as JSON
	err := json.NewDecoder(r.Body).Decode(&bod)

	// Send 400 is error occurs
	if err != nil {
		SendJSON(w, Response{Status: http.StatusBadRequest, State: "fail", Result: "error: invalid input"})
		return
	}
	// Make sure name is provided
	if bod.Name == "" {
		SendJSON(w, Response{Status: http.StatusBadRequest, State: "fail", Result: "error: must provide a name"})
		return
	}
	// Make sure url is provided
	if bod.URL == "" {
		SendJSON(w, Response{Status: http.StatusBadRequest, State: "fail", Result: "error: must provide a url"})
		return
	}

	// If slug is not provided, generate random one
	if bod.Slug == "" {
		slug = CreateSlug(5)
	} else {
		slug = bod.Slug
	}

	// Find duplicate
	res, cancel, _ := FindSlug(bson.M{"slug": slug})
	defer cancel()

	if res["slug"] == slug {
		SendJSON(w, Response{Status: http.StatusConflict, State: "fail", Result: fmt.Sprintf("slug with identifier '%s' already exists", slug)})
		return
	}

	// Create url in database
	timestamp := time.Now().String()
	_, err = db.InsertOne(ctx, bson.M{
		"slug":       slug,
		"name":       bod.Name,
		"url":        bod.URL,
		"created-at": timestamp,
	})

	if err != nil {
		// Handle duplicate error
		if strings.IndexAny(err.Error(), "E11000 duplicate key error collection") != -1 {
			SendJSON(w, Response{Status: http.StatusConflict, State: "fail", Result: fmt.Sprintf("slug with identifier '%s' already exists", slug)})
			return
		}
		panic(err)
	}

	SendJSON(w, Response{Status: http.StatusOK, State: "ok", Result: Tinyer{Slug: slug, Name: bod.Name, CreatedAt: timestamp}})
}

// DeleteURL will delete a url
func DeleteURL(w http.ResponseWriter, r *http.Request) {
	id := s.Make(mux.Vars(r)["id"])
	res, cancel, _ := FindSlug(bson.M{"slug": id})
	defer cancel()

	// No item found
	if res == nil {
		SendJSON(w, Response{Status: http.StatusNotFound, State: "fail", Result: fmt.Sprintf("error: url with identifier '%s' could not be found", id)})
		return
	}

	// Delete url if found
	if res["slug"] == id {
		if _, err := db.DeleteOne(ctx, bson.M{"slug": id}); err != nil {
			panic(err)
		}

		SendJSON(w, Response{Status: http.StatusOK, State: "ok", Result: res})
		return
	}

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

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

// CreateSlug will create a slug with random charas given the provided length
func CreateSlug(n int) string {
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return s.Make(*(*string)(unsafe.Pointer(&b)))
}

// FindSlug finds item in database with same slug as provided
// Yes I know I could do 'FindOne' however when I tried with that
// The codebase wouldn't detect the database updating and would error
// This seems to prevent that so ¯\_(ツ)_/¯
func FindSlug(filter bson.M) (bson.M, context.CancelFunc, error) {
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
	// Find filter in db
	cur, _ := db.Find(ctx, filter)

	// Item could not be found
	if cur.RemainingBatchLength() <= 0 {
		return nil, cancel, errors.New("no items found")
	}

	var item bson.M

	for cur.Next(ctx) {
		var res bson.M
		cur.Decode(&res)

		if res != nil {
			item = res
			break
		}
	}

	if err := cur.Err(); err != nil {
		panic(err)
	}

	return item, cancel, nil
}
