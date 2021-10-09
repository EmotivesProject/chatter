package main

import (
	"chatter/internal/api"
	"chatter/internal/connections"
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/middlewares"
	commonPostgres "github.com/TomBowyerResearchProject/common/postgres"
	"github.com/TomBowyerResearchProject/common/verification"
)

const timeBeforeTimeout = 15

func main() {
	initServices()

	go connections.HandleMessages()

	router := api.CreateRouter()

	srv := http.Server{
		Handler:      router,
		Addr:         os.Getenv("HOST") + ":" + os.Getenv("PORT"),
		WriteTimeout: timeBeforeTimeout * time.Second,
		ReadTimeout:  timeBeforeTimeout * time.Second,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		logger.Infof("Shutting down server")

		if err := srv.Shutdown(context.Background()); err != nil {
			logger.Infof("HTTP server Shutdown: %v", err)
		}

		commonPostgres.CloseDatabase()

		logger.Infof("mongo disconnected")
		close(idleConnsClosed)
	}()

	logger.Info("Starting Server")

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		logger.Infof("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

func initServices() {
	logger.InitLogger("chatter", logger.EmailConfig{
		From:     os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		Level:    os.Getenv("EMAIL_LEVEL"),
	})

	verification.Init(verification.VerificationConfig{
		VerificationURL:     os.Getenv("VERIFICATION_URL"),
		AuthorizationSecret: os.Getenv("NOTIFICATION_AUTH"),
	})

	middlewares.Init(middlewares.Config{
		AllowedOrigins: os.Getenv("ALLOWED_ORIGINS"),
		AllowedMethods: "GET,POST,OPTIONS,DELETE",
		AllowedHeaders: "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, X-CSRF-Token",
	})

	err := commonPostgres.Connect(commonPostgres.Config{
		URI: "postgres://postgres@postgres_db:5432/chatter_db",
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
