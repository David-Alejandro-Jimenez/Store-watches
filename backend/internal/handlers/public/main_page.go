package public

import (
	"net/http"
	"os"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	var filepath = "./../frontend/index.html"
	var _, err = os.Stat(filepath)
	if err != nil {
		http.Error(w, "File not found", http.StatusInternalServerError)
	}

	http.ServeFile(w, r, filepath)
}
