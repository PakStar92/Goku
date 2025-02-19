// index.go
package main

import (
	"log"
	"net/http"
	"github.com/PakStar92/Goku//routes"  // Import the routes package
)

func main() {
	// Initialize API and other routes
	routes.InitRoutes()

	// Start the server
	log.Println("Server started at :3000")
	log.Fatal(http.ListenAndServe(":3000", nil)
           )
}
