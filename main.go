package main

import (
	"context"
	"errors"
	"fmt"
	"main/constants"
	startupModels "main/models/startup"
	"main/startup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	consulClient "github.com/punk-link/consul-client"
	envManager "github.com/punk-link/environment-variable-manager"
	"github.com/punk-link/logger"
)

func main() {
	logger := logger.New()

	environmentName := getEnvironmentName()
	logger.LogInfo("Artist Updater API is running as '%s'", environmentName)

	consul, err := getConsulClient(constants.SERVICE_NAME, environmentName)
	if err != nil {
		logger.LogFatal(err, "Can't initialize the consul client: '%s'", err.Error())
		return
	}

	hostSettingsValues, err := consul.Get("HostSettings")
	if err != nil {
		logger.LogFatal(err, "Can't obtain host settings from Consul: '%s'", err.Error())
		return
	}
	hostSettings := hostSettingsValues.(map[string]any)

	app := startup.Configure(logger, consul, &startupModels.StartupOptions{
		EnvironmentName: environmentName,
		GinMode:         hostSettings["Mode"].(string),
		ServiceName:     constants.SERVICE_NAME,
	})
	app.Run()

	hostAddress := hostSettings["Address"]
	hostPort := hostSettings["Port"]
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", hostAddress, hostPort),
		Handler: app,
	}

	go func() {
		logger.LogInfo("Starting...")
		err := server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			logger.LogError(err, "Listen error: %s\n", err.Error())
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

func getConsulClient(storageName string, environmentName string) (*consulClient.ConsulClient, error) {
	isExist, consulAddress := envManager.TryGetEnvironmentVariable("PNKL_CONSUL_ADDR")
	if !isExist {
		return nil, fmt.Errorf("can't find value of the '%s' environment variable", "PNKL_CONSUL_ADDR")
	}

	isExist, consulToken := envManager.TryGetEnvironmentVariable("PNKL_CONSUL_TOKEN")
	if !isExist {
		return nil, fmt.Errorf("can't find value of the '%s' environment variable", "PNKL_CONSUL_TOKEN")
	}

	consul, err := consulClient.New(&consulClient.ConsulConfig{
		Address:         consulAddress,
		EnvironmentName: environmentName,
		StorageName:     storageName,
		Token:           consulToken,
	})

	return consul, err
}

func getEnvironmentName() string {
	isExist, name := envManager.TryGetEnvironmentVariable("GO_ENVIRONMENT")
	if !isExist {
		return "Development"
	}

	return name
}
