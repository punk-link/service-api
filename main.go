package main

import (
	"context"
	"errors"
	"fmt"
	"main/data"
	"main/infrastructure"
	"main/startup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/rs/zerolog/log"
)

func main() {
	environmentName := "Development"

	envName := infrastructure.GetEnvironmentVariable("GO_ENVIRONMENT")
	if envName != "" {
		environmentName = envName
	}

	log.Info().Msgf("Artist Updater API is running as '%s'", environmentName)

	client, err := api.NewClient(&api.Config{
		Address: "65.21.241.207:8500",
		Scheme:  "http",
		Token:   "12a106ef-0df4-0536-3cac-d21f3877cffd",
	})
	if err != nil {
		log.Error().Msgf("Listen error: %s\n", err)
		panic(err)
	}

	kv := client.KV()
	pair, _, err := kv.Get("/development/connection", nil)
	if err != nil {
		panic(err)
	}

	log.Info().Msgf("KV: %v %s\n", pair.Key, pair.Value)

	hostAddress := infrastructure.GetEnvironmentVariable("GIN_ADDR")
	hostPort := infrastructure.GetEnvironmentVariable("GIN_PORT")

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
