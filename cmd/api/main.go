package main

import (
	"chatter/internal/api"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	router := api.CreateRouter()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(host+":"+port, router))
}
