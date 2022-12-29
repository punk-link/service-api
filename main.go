package main

import (
	"errors"
	"fmt"
	"main/constants"
	presentationGrpc "main/grpc/presentation"
	startupModels "main/models/startup"
	"main/startup"
	"net"

	consulClient "github.com/punk-link/consul-client"
	envManager "github.com/punk-link/environment-variable-manager"
	"github.com/punk-link/logger"
	vaultClient "github.com/punk-link/vault-client"
	"google.golang.org/grpc"
)

func main() {
	logger := logger.New()
	envManager := envManager.New()

	environmentName := getEnvironmentName(envManager)
	logger.LogInfo("Artist Updater API is running as '%s'", environmentName)

	appSecrets := getSecrets(envManager, logger, SECRET_ENGINE_NAME, constants.SERVICE_NAME)
	consul := getConsulClient(envManager, logger, appSecrets, constants.SERVICE_NAME, environmentName)

	go runGrpc(logger, consul)

	runHttp(logger, consul, appSecrets, environmentName)
}

func getConsulClient(envManager envManager.EnvironmentVariableManager, logger logger.Logger, appSecrets map[string]any, storageName string, environmentName string) consulClient.ConsulClient {
	client, err := consulClient.New(&consulClient.ConsulConfig{
		Address:         appSecrets["consul-address"].(string),
		EnvironmentName: environmentName,
		StorageName:     storageName,
		Token:           appSecrets["consul-token"].(string),
	})

	if err != nil {
		logger.LogFatal(err, "Can't initialize Consul client: '%s'", err.Error())
	}

	return client
}

func getSecrets(envManager envManager.EnvironmentVariableManager, logger logger.Logger, storeName string, secretName string) map[string]any {
	vaultAddress, isExist := envManager.TryGet("PNKL_VAULT_ADDR")
	if !isExist {
		err := errors.New("can't get PNKL_VAULT_ADDR environment variable")
		logger.LogFatal(err, err.Error())
	}

	vaultToken, isExist := envManager.TryGet("PNKL_VAULT_TOKEN")
	if !isExist {
		err := errors.New("an't get PNKL_VAULT_TOKEN environment variable")
		logger.LogFatal(err, err.Error())
	}

	vaultConfig := &vaultClient.VaultClientOptions{
		Endpoint: vaultAddress,
		RoleName: secretName,
	}

	vaultClient := vaultClient.New(vaultConfig, logger)
	return vaultClient.Get(vaultToken, storeName, secretName)
}

func getEnvironmentName(envManager envManager.EnvironmentVariableManager) string {
	name, isExist := envManager.TryGet("GO_ENVIRONMENT")
	if !isExist {
		return "Development"
	}

	return name
}

func runGrpc(logger logger.Logger, consul consulClient.ConsulClient) {
	grpcSettingsValues, err := consul.Get("GrpcSettings")
	if err != nil {
		logger.LogFatal(err, "Can't obtain host settings from Consul: '%s'", err.Error())
	}
	grpcSettings := grpcSettingsValues.(map[string]any)
	grpcPresentationPort := grpcSettings["PresentationPort"]

	presentationPortListener, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPresentationPort))
	if err != nil {
		logger.LogFatal(err, "GRPC client failed to listen: %v", err)
	}

	logger.LogInfo("Listening and serving TCP on :%s", grpcPresentationPort)

	presentationServer := grpc.NewServer()
	presentationGrpc.RegisterPresentationServer(presentationServer, &presentationGrpc.Server{})
	if err := presentationServer.Serve(presentationPortListener); err != nil {
		logger.LogFatal(err, "Failed to serve: %v", err)
	}
}

func runHttp(logger logger.Logger, consul consulClient.ConsulClient, appSecrets map[string]any, environmentName string) {
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
	app.Run(fmt.Sprintf(":%s", hostSettings["Port"]))

	// hostAddress := hostSettings["Address"]
	// hostPort := hostSettings["Port"]
	// server := &http.Server{
	// 	Addr:    fmt.Sprintf("%s:%s", hostAddress, hostPort),
	// 	Handler: app,
	// }

	// go func() {
	// 	logger.LogInfo("Starting...")
	// 	err := server.ListenAndServe()
	// 	if err != nil && errors.Is(err, http.ErrServerClosed) {
	// 		logger.LogError(err, "Listen error: %s\n", err.Error())
	// 	}
	// }()

	// quit := make(chan os.Signal, 1)
	// signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// <-quit

	// logger.LogInfo("Sutting down...")

	// ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancelFunc()

	// if err := server.Shutdown(ctx); err != nil {
	// 	logger.LogError(err, "Server forced to shutdown: %s", err)
	// }

	// logger.LogInfo("Exiting")
}

const SECRET_ENGINE_NAME = "secrets"
