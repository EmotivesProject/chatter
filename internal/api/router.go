package api

import (
	"net/http"

	"github.com/TomBowyerResearchProject/common/middlewares"
	"github.com/TomBowyerResearchProject/common/response"
	"github.com/TomBowyerResearchProject/common/verification"
	"github.com/go-chi/chi"
)

func CreateRouter() http.Handler {
	r := chi.NewRouter()

	r.Use(middlewares.SimpleMiddleware())

	r.Get("/healthz", response.Healthz)

	r.With(verification.VerifyToken()).Route("/user", func(r chi.Router) {
		r.Post("/", createUser)
	})

	r.With(verification.VerifyJTW()).Route("/ws_token", func(r chi.Router) {
		r.Get("/", createTocken)
	})

	r.With(verification.VerifyJTW()).Route("/messages", func(r chi.Router) {
		r.Get("/", getMessages)
	})

	r.Route("/connections", func(r chi.Router) {
		r.Get("/", getConnectedUsers)
	})

	r.Route("/ws", func(r chi.Router) {
		r.Get("/", handleConnections)
	})

	return r
}
