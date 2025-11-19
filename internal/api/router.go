package api

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Register(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(3 * time.Second))

	// CORS wie in Hono
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET"},
		AllowedHeaders: []string{"*"},
		MaxAge:         300,
	}))

	r.Get("/health", handleHealth)
	// TODO all errors should be logged
	r.Post("/webhook", handleWebHook)

	r.Route("/api", func(api chi.Router) {
		api.Get("/*", handleContent)
	})
}
