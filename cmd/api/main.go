package main

import (
	"log"
	"net/http"

	server "dynoc-registry/internal/http"
)

func main() {
	srv := server.NewServer()

	log.Println("DynoC registry API listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}
