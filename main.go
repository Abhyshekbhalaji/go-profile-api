package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Profile struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Posts    int    `json:"posts"`
	Reels    int    `json:"reels"`
}

var (
	usersSampleDb = map[string]Profile{
		"1": {
			Name:     "Abhyshek",
			Username: "nover._4._u",
			Posts:    2,
			Reels:    0,
		},
	}
	nextID    = 2
	mutexLock sync.RWMutex
)

func main() {
	log.Fatal(http.ListenAndServe(":8080", newMux()))
}

// newMux builds the application's router so both main and HTTP tests use
// the same route configuration.
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /users/{userId}", GETHandler)
	mux.HandleFunc("POST /users", POSTHandler)
	mux.HandleFunc("PUT /users/{userId}", PUTHandler)
	mux.HandleFunc("DELETE /users/{userId}", DELETEHandler)
	return mux
}

func GETHandler(w http.ResponseWriter, r *http.Request) {
	profile, exists := Dbcall(r.PathValue("userId"))
	if !exists {
		http.Error(w, "No such profile exists", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		fmt.Printf("Unable to encode profile: %v\n", err)
	}
}

// Dbcall returns a copy of a profile and whether the ID exists.
func Dbcall(userID string) (Profile, bool) {
	if userID == "" {
		return Profile{}, false
	}

	mutexLock.RLock()
	defer mutexLock.RUnlock()

	profile, exists := usersSampleDb[userID]
	return profile, exists
}

func POSTHandler(w http.ResponseWriter, r *http.Request) {
	var profile Profile

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&profile); err != nil {
		http.Error(w, "Unable to decode the data", http.StatusBadRequest)
		return
	}
	if !validProfile(profile) {
		http.Error(w, "Invalid profile data", http.StatusBadRequest)
		return
	}

	mutexLock.Lock()
	id := strconv.Itoa(nextID)
	usersSampleDb[id] = profile
	nextID++
	mutexLock.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "/users/"+id)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		fmt.Printf("Unable to encode profile: %v\n", err)
	}
}

func PUTHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	var profile Profile

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&profile); err != nil {
		http.Error(w, "Unable to decode the data", http.StatusBadRequest)
		return
	}
	if !validProfile(profile) {
		http.Error(w, "Invalid profile data", http.StatusBadRequest)
		return
	}

	mutexLock.Lock()
	if _, exists := usersSampleDb[userID]; !exists {
		mutexLock.Unlock()
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	usersSampleDb[userID] = profile
	mutexLock.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(profile); err != nil {
		fmt.Printf("Unable to encode profile: %v\n", err)
	}
}

func DELETEHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	mutexLock.Lock()
	_, exists := usersSampleDb[userID]
	if exists {
		delete(usersSampleDb, userID)
	}
	mutexLock.Unlock()

	if !exists {
		http.Error(w, "No such profile exists", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func validProfile(profile Profile) bool {
	return profile.Name != "" &&
		profile.Username != "" &&
		profile.Posts >= 0 &&
		profile.Reels >= 0
}
