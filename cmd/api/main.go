package main

import (
	"chatter/internal/api"
	"chatter/internal/connections"
	"chatter/internal/db"
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
	commonMongo "github.com/TomBowyerResearchProject/common/mongo"
	"github.com/TomBowyerResearchProject/common/verification"
)

const timeBeforeTimeout = 15

func main() {
	logger.InitLogger("chatter")

	verification.Init(verification.VerificationConfig{
		VerificationURL: "http://uacl/authorize",
	})

	middlewares.Init(middlewares.Config{
		AllowedOrigin:  "*",
		AllowedMethods: "GET,OPTIONS",
		AllowedHeaders: "*",
	})

	err := commonMongo.Connect(commonMongo.Config{
		URI:    "mongodb://admin:admin@mongo_db:27017",
		DBName: db.DBName,
	})
	if err != nil {
		log.Fatal(err.Error())
	}

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

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			logger.Infof("HTTP server Shutdown: %v", err)
		}

		monogodb := commonMongo.GetDatabase()
		if monogodb != nil {
			_ = monogodb.Client().Disconnect(context.Background())
		}

		logger.Infof("mongo disconnected")

		close(idleConnsClosed)
	}()

	logger.Info("Starting Server")

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		// Error starting or closing listener:
		logger.Infof("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
