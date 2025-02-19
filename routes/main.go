package main

import (
	"net/http"
	"path/filepath"

	"github.com/gorilla/mux"
)

// RegisterMainRoutes sets up the main routes for serving HTML files
func RegisterMainRoutes(router *mux.Router) {
	// Base path for HTML files
	basePath := "./views"

	// Route for /docs
	router.HandleFunc("/docs", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(basePath, "docs.html"))
	})

	// Route for /
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(basePath, "index.html"))
	})

	// Route for /about
	router.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(basePath, "about.html"))
	})
}
