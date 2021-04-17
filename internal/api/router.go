package api

import (
	"github.com/TomBowyerResearchProject/common/verification"
	"github.com/go-chi/chi"
)

func CreateRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(SimpleMiddleware())

	r.Route("/", func(r chi.Router) {
		r.Get("/healthz", healthz)
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
