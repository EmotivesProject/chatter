package main

import (
	"chatter/internal/api"
	"chatter/internal/connections"
	"chatter/internal/db"
	"chatter/internal/logger"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger("chatter")

	router := api.CreateRouter()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	db.Init()

	go connections.HandleMessages()

	log.Fatal(http.ListenAndServe(host+":"+port, router))
}
