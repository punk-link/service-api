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

	vault "github.com/hashicorp/vault/api"
	consulClient "github.com/punk-link/consul-client"
	envManager "github.com/punk-link/environment-variable-manager"
	"github.com/punk-link/logger"
)

func main() {
	logger := logger.New()

	environmentName := getEnvironmentName()
	logger.LogInfo("Artist Updater API is running as '%s'", environmentName)

	appSecrets, err := getSecrets(SECRET_STORAGE_NAME, SERVICE_NAME)
	if err != nil {
		logger.LogFatal(err, "Vault access error: %s", err)
	}

	consul, err := getConsulClient(appSecrets, SERVICE_NAME, environmentName)
	if err != nil {
		logger.LogFatal(err, "Can't initialize Consul client: '%s'", err.Error())
	}

	hostSettingsValues, err := consul.Get("HostSettings")
	if err != nil {
		logger.LogFatal(err, "Can't obtain host settings from Consul: '%s'", err.Error())
	}
	hostSettings := hostSettingsValues.(map[string]any)

	app := startup.Configure(logger, consul, appSecrets, &startupModels.StartupOptions{
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

func getConsulClient(appSecrets map[string]any, storageName string, environmentName string) (*consulClient.ConsulClient, error) {
	return consulClient.New(&consulClient.ConsulConfig{
		Address:         appSecrets["consul-address"].(string),
		EnvironmentName: environmentName,
		StorageName:     storageName,
		Token:           appSecrets["consul-token"].(string),
	})
}

func getSecrets(storeName string, secretName string) (map[string]any, error) {
	var err error = nil
	// isExist, vaultAddress := envManager.TryGetEnvironmentVariable("PNKL_VAULT_ADDR")
	// if !isExist {
	// 	return nil, fmt.Errorf("can't get PNKL_VAULT_ADDR environment variable")
	// }
	vaultAddress := "http://vault.dev.svc.cluster.local:8200"

	isExist, vaultToken := envManager.TryGetEnvironmentVariable("PNKL_VAULT_TOKEN")
	if !isExist {
		return nil, fmt.Errorf("can't get PNKL_VAULT_TOKEN environment variable")
	}

	vaultConfig := vault.DefaultConfig()
	vaultConfig.Address = vaultAddress

	vaultClient, err := vault.NewClient(vaultConfig)
	if err != nil {
		return nil, fmt.Errorf("can't initialize Vault: %s", err)
	}

	vaultClient.SetToken(vaultToken)

	vaultCtx := context.Background()
	serviceApiSettings, err := vaultClient.KVv2(storeName).Get(vaultCtx, secretName)
	if err != nil {
		return nil, fmt.Errorf("can't obtain app secrets: %s", err)
	}

	return serviceApiSettings.Data, nil
}

func getEnvironmentName() string {
	isExist, name := envManager.TryGetEnvironmentVariable("GO_ENVIRONMENT")
	if !isExist {
		return "Development"
	}

	return name
}

const SECRET_STORAGE_NAME = "secrets"
const SERVICE_NAME = "service-api"
