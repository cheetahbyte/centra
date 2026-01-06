package logger

import (
	"io"
	"os"
	"time"

	"github.com/cheetahbyte/centra/internal/config"
	"github.com/rs/zerolog"
)

func AcquireLogger() zerolog.Logger {
	conf := config.Get()
	level, err := zerolog.ParseLevel(conf.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	var out io.Writer
	out = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	if conf.LogStructured {
		out = os.Stdout
	}
	log := zerolog.New(out).With().Timestamp().Logger()
	return log
}
