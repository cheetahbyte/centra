package middleware

import (
	"crypto/subtle"
	"net/http"
	"strings"
	"time"

	"github.com/cheetahbyte/centra/internal/config"
	"github.com/cheetahbyte/centra/internal/logger"
)

func APIKeyAuth() func(http.Handler) http.Handler {
	expectedKey := config.GetAPIKey()
	log := logger.AcquireLogger()

	if expectedKey == "" {
		log.Warn().Msg("no api key configured. api key auth is DISABLED.")
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			authHeader := r.Header.Get("Authorization")
			var providedKey string
			if strings.HasPrefix(authHeader, "Bearer ") {
				providedKey = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				providedKey = r.Header.Get("X-API-Key")
			}

			if providedKey == "" {
				http.Error(w, "missing api key", http.StatusUnauthorized)
				log.Warn().
					Str("path", r.URL.Path).
					Str("method", r.Method).
					Msg("request without API key")
				return
			}

			if subtle.ConstantTimeCompare([]byte(providedKey), []byte(expectedKey)) != 1 {
				http.Error(w, "invalid api key", http.StatusUnauthorized)
				log.Warn().
					Str("path", r.URL.Path).
					Str("method", r.Method).
					Dur("duration", time.Since(start)).
					Msg("request with INVALID API key")
				return
			}

			log.Debug().
				Str("path", r.URL.Path).
				Str("method", r.Method).
				Dur("duration", time.Since(start)).
				Msg("request authenticated with API key")

			next.ServeHTTP(w, r)
		})
	}
}
