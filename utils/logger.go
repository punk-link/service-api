package utils

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	logLevel := strings.ToLower(GetEnvironmentVariable("LOG_LEVEL"))

	switch logLevel {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	}

	host, err := os.Hostname()
	if err != nil {
		log.Logger = log.With().Str("host", "unknown").Logger()
	} else {
		log.Logger = log.With().Str("host", host).Logger()
	}

	log.Logger = log.With().Caller().Logger()
}
