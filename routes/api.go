// routes/api.go
package routes

import (
	"fmt"
	"net/http"
	"github.com/PakStar92/Goku/lib" // Import the lib package from the Goku repo
)

func GreetHandler(w http.ResponseWriter, r *http.Request) {
	// Assume "name" query parameter is passed in the URL
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "Guest"
	}

	greeting := lib.GreetUser(name) // Using a utility function from lib package
	fmt.Fprintf(w, greeting
             )
}
