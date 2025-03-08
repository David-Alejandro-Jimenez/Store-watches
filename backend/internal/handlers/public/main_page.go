package public

import (
	"net/http"
	"os"
)

// The Main_page function is responsible for serving the main HTML file (the home page) of the web application.
// 1. Objective: Serve the main page of the application (index.html file).
// 2. Process:
		// Defines the file path.
		// Verify the existence of the file.
 		// In case of error, an error message is returned.
		// If everything is correct, send the file to the client using http.ServeFile.
// This function is essential to load the main frontend interface in the web application.
func Main_page(w http.ResponseWriter, r *http.Request) {
	var filepath = "./../frontend/index.html"
	var _, err = os.Stat(filepath)
	if err != nil {
		http.Error(w, "File not found", http.StatusInternalServerError)
	}

	http.ServeFile(w, r, filepath)
}