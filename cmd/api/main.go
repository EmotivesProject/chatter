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
)

func main() {
	logger.InitLogger("chatter")

	verification.Init(verification.VerificationConfig{
		VerificationURL: "http://uacl/authorize",
	})

	middlewares.Init(middlewares.Config{
		AllowedOrigin:  "*",
		AllowedMethods: "GET,OPTIONS",
		// nolint:lll
		AllowedHeaders: "Accept, Content-Type, Content-Length, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header",
	})

	router := api.CreateRouter()

	err := commonMongo.Connect(commonMongo.Config{
		URI:    "mongodb://admin:admin@mongo_db:27017",
		DBName: db.DBName,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	go connections.HandleMessages()

	log.Fatal(http.ListenAndServe(os.Getenv("HOST")+":"+os.Getenv("PORT"), router))
}
