package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Creator name
const creator = "Qasim Ali 🦋"

// ErrorMessage represents the structure of an error response
type ErrorMessage struct {
	Status  bool   `json:"status"`
	Creator string `json:"creator"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

// createErrorMessage generates an error message with the given message and code
func createErrorMessage(message string, code int) ErrorMessage {
	return ErrorMessage{
		Status:  false,
		Creator: creator,
		Code:    code,
		Message: message,
	}
}

// Predefined error messages
var loghandler = map[string]ErrorMessage{
	"notparam":    createErrorMessage("Please provide the apikey", 406),
	"noturl":      createErrorMessage("Please provide the url", 406),
	"notquery":    createErrorMessage("Please provide the query", 406),
	"notkata":     createErrorMessage("Please provide the kata", 406),
	"nottext":     createErrorMessage("Please provide the text", 406),
	"nottext2":    createErrorMessage("Please provide the text2", 406),
	"notnabi":     createErrorMessage("Please provide the nabi name", 406),
	"nottext3":    createErrorMessage("Please provide the text3", 406),
	"nottheme":    createErrorMessage("Please provide the theme", 406),
	"notname":     createErrorMessage("Please provide the name", 406),
	"notusername": createErrorMessage("Please provide the username", 406),
	"notvalue":    createErrorMessage("Please provide the value", 406),
	"invalidKey":  createErrorMessage("Invalid apikey", 406),
}

var invalidlink = ErrorMessage{
	Status:  false,
	Creator: creator,
	Message: "Error, the link might be invalid.",
}

var invalidkata = ErrorMessage{
	Status:  false,
	Creator: creator,
	Message: "Error, the word might not exist in the API.",
}

var errorMessage = ErrorMessage{
	Status:  false,
	Creator: creator,
	Message: "An error occurred.",
}

// List of valid API keys
var listkey = []string{"Suhail", "GURU", "APIKEY"}

// RegisterRoutes sets up the API routes
func RegisterRoutes(router *mux.Router) {
	// POST /apikey - Add a new API key
	router.HandleFunc("/apikey", addAPIKey).Methods("POST")

	// DELETE /apikey - Delete an API key
	router.HandleFunc("/apikey", deleteAPIKey).Methods("DELETE")
}

// addAPIKey handles the POST /apikey endpoint
func addAPIKey(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		respondWithError(w, loghandler["notparam"])
		return
	}

	for _, k := range listkey {
		if k == key {
			respondWithJSON(w, http.StatusOK, map[string]string{"message": "API key already registered"})
			return
		}
	}

	listkey = append(listkey, key)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Successfully registered " + key + " in the API key database"})
}

// deleteAPIKey handles the DELETE /apikey endpoint
func deleteAPIKey(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("delete")
	if key == "" {
		respondWithError(w, loghandler["notparam"])
		return
	}

	index := -1
	for i, k := range listkey {
		if k == key {
			index = i
			break
		}
	}

	if index == -1 {
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "API key does not exist"})
		return
	}

	listkey = append(listkey[:index], listkey[index+1:]...)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "API key successfully deleted"})
}

// respondWithError sends an error response in JSON format
func respondWithError(w http.ResponseWriter, err ErrorMessage) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotAcceptable)
	json.NewEncoder(w).Encode(err)
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
