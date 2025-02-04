package main

import (
	"log"
	"net/http"

	"github.com/David-Alejandro-Jimenez/venta-relojes/internal"
	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/config"
)

func startTheServer() {
	var errConfig = config.LoadConfig()
	if errConfig != nil {
		log.Fatalf("Error loading configuration: %v", errConfig)
	}

	var port = ":8080"
	var router = internal.SetupRouter()
	log.Printf("Server listening on http://localhost%s", port)

	var err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func main() {
	startTheServer()
}