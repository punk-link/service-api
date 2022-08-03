package main

import (
	"context"
	"errors"
	"fmt"
	"main/data"
	"main/startup"
	"main/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	environmentName := "Development"

	envName := utils.GetEnvironmentVariable("GO_ENVIRONMENT")
	if envName != "" {
		environmentName = envName
	}

	log.Info().Msgf("Artist Updater API is running as '%s'", environmentName)

	hostAddress := utils.GetEnvironmentVariable("GIN_ADDR")
	hostPort := utils.GetEnvironmentVariable("GIN_PORT")

	app := startup.Configure()
	data.ConfigureDatabase()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", hostAddress, hostPort),
		Handler: app,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Error().Msgf("Listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Sutting down...")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	if err := server.Shutdown(ctx); err != nil {
		log.Error().Msgf("Server forced to shutdown: %s", err)
	}

	log.Info().Msg("Exiting")
}
