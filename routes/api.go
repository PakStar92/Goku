package routes

import (
	"encoding/json"
	"net/http"
	"strings"
	"github.com/gorilla/mux"
	"github.com/rylio/ytdl"
	"os"
	"fmt"
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

// downloadYouTubeVideo handles the YouTube video/audio download
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

	// Use ytdl to fetch video info
	videoInfo, err := ytdl.GetVideoInfo(url)
	if err != nil {
		respondWithError(w, ErrorMessage{
			Status:  false,
			Creator: creator,
			Message: "Failed to fetch YouTube video info",
		})
		return
	}

	// Retrieve available formats
	formats := videoInfo.Formats

	// Select the best audio or video format based on user preference
	format := selectBestFormat(formats)

	// Start downloading the selected format
	if err := downloadFromFormat(format, url, w); err != nil {
		respondWithError(w, ErrorMessage{
			Status:  false,
			Creator: creator,
			Message: "Failed to download the video/audio",
		})
		return
	}

	// Respond to user indicating the download is complete
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"message": "Download complete",
	})
}

// selectBestFormat selects the best video or audio format (you can customize this logic)
func selectBestFormat(formats []ytdl.Format) ytdl.Format {
	var selectedFormat ytdl.Format
	for _, format := range formats {
		// Choose the best format for video (highest quality) or audio (best bitrate)
		if format.HasVideo() && format.HasAudio() {
			if selectedFormat.ID == "" || format.VideoQuality > selectedFormat.VideoQuality {
				selectedFormat = format
			}
		} else if format.HasAudio() {
			if selectedFormat.ID == "" || format.AudioBitrate > selectedFormat.AudioBitrate {
				selectedFormat = format
			}
		}
	}
	return selectedFormat
}

// downloadFromFormat handles the actual download of video/audio in the selected format
func downloadFromFormat(format ytdl.Format, url string, w http.ResponseWriter) error {
	// Create a temp file to store the downloaded content
	tempFile, err := os.CreateTemp("", "downloaded_video_*")
	if err != nil {
		return err
	}
	defer tempFile.Close()

	// Open the video/audio stream based on the selected format
	downloadURL, err := format.DownloadURL(url)
	if err != nil {
		return err
	}

	// Open the stream for downloading
	resp, err := http.Get(downloadURL.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Copy the content to the temp file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		return err
	}

	// Once the file is downloaded, send it as a response to the user
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", tempFile.Name()))
	w.Header().Set("Content-Type", "application/octet-stream")
	http.ServeFile(w, r, tempFile.Name())

	return nil
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

