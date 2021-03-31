package api

import (
	"github.com/go-chi/chi"
)

func CreateRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(SimpleMiddleware())

	r.Route("/", func(r chi.Router) {
		r.Get("/healthz", healthz)
	})

	r.Route("/user", func(r chi.Router) {
		r.Post("/", createUser)
	})

	r.With(verifyJTW()).Route("/ws_token", func(r chi.Router) {
		r.Get("/", createTocken)
	})

	r.Route("/ws", func(r chi.Router) {
		r.Get("/", handleConnections)
	})

	return r
}
