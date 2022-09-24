package common

import (
	"github.com/rs/zerolog/log"
	"github.com/samber/do"
)

type Logger struct{}

func ConstructLogger(injector *do.Injector) (*Logger, error) {
	return ConstructLoggerWithoutInjection(), nil
}

func ConstructLoggerWithoutInjection() *Logger {
	return &Logger{}
}

func (logger *Logger) LogError(err error, format string, args ...interface{}) {
	log.Error().Err(err).Msgf(format, args...)
}

func (logger *Logger) LogFatal(err error, format string, args ...interface{}) {
	log.Fatal().Err(err).Msgf(format, args...)
}

func (logger *Logger) LogInfo(format string, args ...interface{}) {
	log.Info().Msgf(format, args...)
}

func (logger *Logger) LogWarn(format string, args ...interface{}) {
	log.Warn().Msgf(format, args...)
}

func (logger *Logger) Printf(format string, values ...interface{}) {
	log.Printf(format, values...)
}
