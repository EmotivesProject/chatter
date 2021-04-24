package main

import (
	"chatter/internal/api"
	"chatter/internal/connections"
	"chatter/internal/db"
	"log"
	"net/http"
	"os"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/middlewares"
	commonMongo "github.com/TomBowyerResearchProject/common/mongo"
	"github.com/TomBowyerResearchProject/common/verification"
	"github.com/joho/godotenv"
)

func main() {
	logger.InitLogger("chatter")

	verification.Init(verification.VerificationConfig{
		VerificationURL: "http://uacl/authorize",
	})

	middlewares.Init(middlewares.Config{
		AllowedOrigin:  "*",
		AllowedMethods: "GET,OPTIONS",
		AllowedHeaders: "Accept, Content-Type, Content-Length, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header",
	})

	router := api.CreateRouter()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	commonMongo.Connect(commonMongo.Config{
		URI:    "mongodb://admin:admin@mongo:27017",
		DBName: db.DBName,
	})

	go connections.HandleMessages()

	log.Fatal(http.ListenAndServe(host+":"+port, router))
}
