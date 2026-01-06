package helper

import (
	"github.com/cheetahbyte/centra/internal/config"
	"github.com/go-chi/cors"
)

func NewCORSConfig() cors.Options {
	conf := config.Get()

	return cors.Options{
		AllowedOrigins:   conf.AllowedOrigins,
		AllowedMethods:   conf.AllowedMethods,
		AllowedHeaders:   conf.AllowedHeaders,
		AllowCredentials: conf.Credentials,
		ExposedHeaders:   conf.ExposedHeaders,
		MaxAge:           conf.MaxAge,
	}
}
