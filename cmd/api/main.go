package main

import (
	"log"
	"net/http"

	api "dynoc-registry/internal/http"
	jwt "dynoc-registry/internal/jwt"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}
	jwt.Init()

	server := api.NewServer()

	log.Println("DynoC registry API listening on :8080")
	log.Fatal(http.ListenAndServe(":3000", server))
}
