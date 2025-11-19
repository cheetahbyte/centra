package helper

import (
	"github.com/cheetahbyte/centra/internal/config"
	"github.com/go-chi/cors"
)

func NewCORSConfig() cors.Options {
	allowedOrigins := config.GetCorsAllowedOrigins()
	allowedMethods := config.GetCorsAllowedMethods()
	allowedHeaders := config.GetCorsAllowedHeaders()
	allowCredentials := config.GetCorsAllowCredentials()
	exposedHeaders := config.GetCorsExposedHeaders()
	maxAge := config.GetCorsMaxAge()

	return cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		AllowCredentials: allowCredentials,
		ExposedHeaders:   exposedHeaders,
		MaxAge:           maxAge,
	}
}
