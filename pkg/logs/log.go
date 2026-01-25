package logs

import (
	"os"

	"github.com/rs/zerolog"
)

var logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: "2006-01-02 15:04:05 Mon"}).
	Level(zerolog.ErrorLevel).With().Timestamp().Caller().Logger()

func SetLogger(l zerolog.Logger) {
	logger = l
}

func Logger() *zerolog.Logger {
	return &logger
}
