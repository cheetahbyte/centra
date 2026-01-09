package api

import (
	"time"

	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/handlers"
	"github.com/cheetahbyte/centra/internal/helper"
	centraMiddleware "github.com/cheetahbyte/centra/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func Register(r *chi.Mux) {
	conf := config.Get()
	r.Use(middleware.RequestID)
	r.Use(centraMiddleware.LoggingMiddleware(helper.AcquireLogger()))
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(3 * time.Second))

	r.Use(cors.Handler(helper.NewCORSConfig()))

	r.Get("/health", handlers.HandleLivez)
	r.Get("/livez", handlers.HandleLivez)
	r.Get("/readyz", handlers.HandleReadyz)
	r.Post("/webhook", handlers.HandleWebHook)

	r.Route("/api", func(api chi.Router) {
		api.Use(httprate.LimitByIP(conf.RateQuota, time.Minute))

		api.Use(centraMiddleware.APIKeyAuth())
		api.Get("/*", handlers.HandleContent)
	})
}
