package api

import (
	"time"

	"github.com/cheetahbyte/centra/internal/helper"
	"github.com/cheetahbyte/centra/internal/logger"
	centraMiddleware "github.com/cheetahbyte/centra/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Register(r *chi.Mux) {
	r.Use(middleware.RequestID)
	r.Use(centraMiddleware.LoggingMiddleware(logger.AcquireLogger()))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(3 * time.Second))

	r.Use(cors.Handler(helper.NewCORSConfig()))

	r.Get("/health", handleHealth)
	// TODO all errors should be logged
	r.Post("/webhook", handleWebHook)

	r.Route("/api", func(api chi.Router) {
		api.Use(centraMiddleware.APIKeyAuth())
		api.Get("/*", handleContent)
	})
}
