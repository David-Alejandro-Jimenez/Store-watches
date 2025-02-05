package main

import (
	"log"
	"net/http"

	"github.com/David-Alejandro-Jimenez/venta-relojes/internal"
	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/config"
	"github.com/David-Alejandro-Jimenez/venta-relojes/internal/repository/database"
)

func loadConfiguration() {
	var errConfig = config.LoadConfig()
	if errConfig != nil {
		log.Fatalf("Error loading configuration: %v", errConfig)
	}

}

func startDatabase() {
	var errdb = database.InitDB()
	if errdb != nil {
		log.Println("Did not connect to the database")
	}
	defer database.DB.Close()
}

func startTheServer() {
	var port = ":8080"
	var router = internal.SetupRouter()
	log.Printf("Server listening on http://localhost%s", port)

	var err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func main() {
	loadConfiguration()
	startDatabase()
	startTheServer()
}