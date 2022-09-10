package main

import (
	"context"
	"errors"
	"fmt"
	"main/infrastructure"
	"main/infrastructure/consul"
	"main/services/common"
	"main/startup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := common.ConstructLogger()

	environmentName := infrastructure.GetEnvironmentName()
	logger.LogInfo("Artist Updater API is running as '%s'", environmentName)

	consul := consul.BuildConsulClient(logger, "service-api")

	hostSettings := consul.Get("HostSettings").(map[string]interface{})
	hostAddress := hostSettings["Address"]
	hostPort := hostSettings["Port"]
	ginMode := hostSettings["Mode"].(string)

	app := startup.Configure(logger, consul, ginMode)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", hostAddress, hostPort),
		Handler: app,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.LogError(err, "Listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.LogInfo("Sutting down...")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	if err := server.Shutdown(ctx); err != nil {
		logger.LogError(err, "Server forced to shutdown: %s", err)
	}

	logger.LogInfo("Exiting")
}
