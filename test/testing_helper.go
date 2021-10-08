package test

import (
	"chatter/internal/api"
	"log"
	"math/rand"
	"net/http/httptest"
	"os"
	"time"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/middlewares"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
	"github.com/TomBowyerResearchProject/common/verification"
)

var TS *httptest.Server

func SetUpIntegrationTest() {
	rand.Seed(time.Now().Unix())

	logger.InitLogger("chatter", logger.EmailConfig{
		From:     os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Level:    os.Getenv("EMAIL_LEVEL"),
	})

	middlewares.Init(middlewares.Config{
		AllowedOrigins: "*",
		AllowedMethods: "GET,POST,OPTIONS,DELETE",
		AllowedHeaders: "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token",
	})

	verification.Init(verification.VerificationConfig{
		VerificationURL: "http://0.0.0.0:8082/authorize",
	})

	err := commonPostgres.Connect(commonPostgres.Config{
		URI: "postgres://tom:tom123@localhost:5435/chatter_db",
	})
	if err != nil {
		log.Fatal(err.Error())
	}

	router := api.CreateRouter()

	TS = httptest.NewServer(router)
}

func TearDownIntegrationTest() {
	TS.Close()
}
