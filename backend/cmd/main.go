package main

import (
	"log"
	"net/http"

	"github.com/David-Alejandro-Jimenez/venta-relojes/internal"
)

func startTheServer() {
	var port = ":8080"
	var router = internal.SetupRouter()
	log.Printf("Servidor escuchando en http://localhost%s", port)

	var err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}

func main() {
	startTheServer()
}