package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/skip2/go-qrcode" // For QR code generation
	"github.com/kkdai/youtube/v2" // For YouTube download
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

// RegisterAPIRoutes sets up the API routes
func RegisterAPIRoutes(router *mux.Router) {
	// GET /apikey - Check if API key is valid
	router.HandleFunc("/apikey", getAPIKey).Methods("GET")

	// POST /apikey - Add a new API key
	router.HandleFunc("/apikey", addAPIKey).Methods("POST")

	// DELETE /apikey - Delete an API key
	router.HandleFunc("/apikey", deleteAPIKey).Methods("DELETE")

	// GET /qrcode - Generate QR Code
	router.HandleFunc("/qrcode", generateQRCode).Methods("GET")

	// GET /ytdl - Download YouTube video info
	router.HandleFunc("/ytdl", downloadYouTubeVideo).Methods("GET")
}

// getAPIKey handles the GET /apikey endpoint
func getAPIKey(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	if key == "" {
		respondWithError(w, loghandler["notparam"])
		return
	}

	if !isValidAPIKey(key) {
		respondWithError(w, loghandler["invalidKey"])
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":  true,
		"creator": creator,
		"message": "API key is valid",
	})
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

// generateQRCode handles the GET /qrcode endpoint
func generateQRCode(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	apiKey := r.URL.Query().Get("apikey")
	if !isValidAPIKey(apiKey) {
		respondWithError(w, loghandler["invalidKey"])
		return
	}

	// Get text from query parameters
	text := r.URL.Query().Get("text")
	if text == "" {
		respondWithError(w, loghandler["nottext"])
		return
	}

	// Generate QR code
	qrCode, err := qrcode.Encode(text, qrcode.Medium, 256)
	if err != nil {
		respondWithError(w, ErrorMessage{
			Status:  false,
			Creator: creator,
			Message: "Failed to generate QR code",
		})
		return
	}

	// Return QR code as base64-encoded JSON
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":  true,
		"creator": creator,
		"qrcode":  qrCode,
	})
}

func downloadYouTubeVideo(w http.ResponseWriter, r *http.Request) {
	// Validate API key
	apiKey := r.URL.Query().Get("apikey")
	if !isValidAPIKey(apiKey) {
		respondWithError(w, loghandler["invalidKey"])
		return
	}

	// Get YouTube URL from query parameters
	url := r.URL.Query().Get("url")
	if url == "" {
		respondWithError(w, loghandler["noturl"])
		return
	}

	// Validate YouTube URL
	if !strings.Contains(url, "youtube.com") && !strings.Contains(url, "youtu.be") {
		respondWithError(w, invalidlink)
		return
	}

	// Download YouTube video info
	client := youtube.Client{}
	video, err := client.GetVideo(url)
	if err != nil {
		respondWithError(w, ErrorMessage{
			Status:  false,
			Creator: creator,
			Message: "Failed to fetch YouTube video info",
		})
		return
	}

	// Extract available formats (video and audio)
	formats := make([]map[string]interface{}, 0)
	for _, format := range video.Formats {
		formats = append(formats, map[string]interface{}{
			"itag":          format.ItagNo,
			"url":           format.URL,
			"mimeType":      format.MimeType,
			"quality":       format.Quality,
			"qualityLabel":  format.QualityLabel,
			"audioChannels": format.AudioChannels,
			"bitrate":       format.Bitrate,
			"fps":           format.FPS,
			"width":         format.Width,
			"height":        format.Height,
		})
	}

	// Extract thumbnails
	thumbnails := make([]map[string]interface{}, 0)
	for _, thumbnail := range video.Thumbnails {
		thumbnails = append(thumbnails, map[string]interface{}{
			"url":    thumbnail.URL,
			"width":  thumbnail.Width,
			"height": thumbnail.Height,
		})
	}

	// Respond with all video info
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":      true,
		"creator":     creator,
		"title":       video.Title,
		"author":      video.Author,
		"description": video.Description,
		"length":      video.Duration.String(),
		"viewCount":   video.ViewCount,
		"publishDate": video.PublishDate,
		"formats":     formats,
		"thumbnails":  thumbnails,
	})
}

// isValidAPIKey checks if the provided API key is valid
func isValidAPIKey(apiKey string) bool {
	for _, key := range listkey {
		if key == apiKey {
			return true
		}
	}
	return false
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
