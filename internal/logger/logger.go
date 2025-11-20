package logger

import (
	"io"
	"os"
	"time"

	"github.com/cheetahbyte/centra/internal/config"
	"github.com/rs/zerolog"
)

func AcquireLogger() zerolog.Logger {
	logLevel := config.GetLogLevel()
	zerolog.SetGlobalLevel(logLevel)
	structedLogging := config.GetLogStructured()

	var out io.Writer
	out = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	if structedLogging {
		out = os.Stdout
	}
	log := zerolog.New(out).With().Timestamp().Logger()
	return log
}
