package api

import (
	"time"

	apihandlers "github.com/cheetahbyte/centra/internal/api/api_handlers"
	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/helper"
	"github.com/cheetahbyte/centra/internal/logger"
	centraMiddleware "github.com/cheetahbyte/centra/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func Register(r *chi.Mux) {
	conf := config.Get()
	r.Use(middleware.RequestID)
	r.Use(centraMiddleware.LoggingMiddleware(logger.AcquireLogger()))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(3 * time.Second))

	r.Use(cors.Handler(helper.NewCORSConfig()))

	r.Get("/health", handleHealth)
	r.Post("/webhook", apihandlers.HandleWebHook)

	r.Route("/api", func(api chi.Router) {
		api.Use(httprate.LimitByIP(conf.RateQuota, time.Minute))

		api.Use(centraMiddleware.APIKeyAuth())
		api.Get("/*", handleContent)
	})
}
