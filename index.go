package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/PakStar92/Goku/routes" // Import the routes package
	"github.com/rs/cors"
)

func main() {
	// Load environment variables (if using .env, use a package like "github.com/joho/godotenv")
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port
	}

	// Initialize the router
	router := mux.NewRouter()

	// Middleware
	router.Use(cors.Default().Handler) // CORS middleware
	// SSL middleware is not directly available in Go; you would typically handle this at the reverse proxy level (e.g., Nginx, Caddy)

	// Serve static files
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	// Initialize routes
	routes.InitRoutes(router)

	// Start the server
	log.Printf("Server started at :%s\n", port)

	// KeepAlive functionality
	go keepAlive()

	// Start the HTTP server
	log.Fatal(http.ListenAndServe(":"+port, router))
}

// keepAlive periodically pings a URL to keep the server alive
func keepAlive() {
	url := os.Getenv("APP_URL")
	if url == "" {
		log.Println("No APP_URL provided, skipping keepAlive...")
		return
	}

	// Validate URL format (basic check)
	if !isValidURL(url) {
		log.Println("Invalid APP_URL format, skipping keepAlive...")
		return
	}

	// Periodically ping the URL every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Error pinging URL: %v\n", err)
			continue
		}
		resp.Body.Close()
	}
}

// isValidURL performs a basic URL format check
func isValidURL(url string) bool {
	return strings.HasPrefix(url, "http://") || strings.HasPrefix(url, "https://")
}
