package logger

import (
	"github.com/punk-link/logger"
	"github.com/samber/do"
)

func New(injector *do.Injector) (logger.Logger, error) {
	return NewWithoutInjection(), nil
}

func NewWithoutInjection() logger.Logger {
	return logger.New()
}
