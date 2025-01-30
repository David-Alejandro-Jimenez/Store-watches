package handlers

import (
	"net/http"
	"os"
)

func Main_page(w http.ResponseWriter, r *http.Request) {
	var filepath = "./../frontend/index.html"
	var _, err = os.Stat(filepath)
	if err != nil {
		http.Error(w, "File not found", http.StatusInternalServerError)
	}

	http.ServeFile(w, r, filepath)

}