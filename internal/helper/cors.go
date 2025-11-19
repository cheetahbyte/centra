package helper

import (
	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/domain"
)

func NewCORSConfig() *domain.CORSConfig {
	allowedOrigins := config.GetCorsAllowedOrigins()
	allowedMethods := config.GetCorsAllowedMethods()
	allowedHeaders := config.GetCorsAllowedHeaders()
	allowCredentials := config.GetCorsAllowCredentials()
	exposedHeaders := config.GetCorsExposedHeaders()
	maxAge := config.GetCorsMaxAge()

	return &domain.CORSConfig{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   allowedMethods,
		AllowedHeaders:   allowedHeaders,
		AllowCredentials: allowCredentials,
		ExposedHeaders:   exposedHeaders,
		MaxAge:           maxAge,
	}
}
