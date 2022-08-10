package common

import "github.com/rs/zerolog/log"

type Logger struct{}

func BuildLogger() *Logger {
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
