package consul

import (
	"encoding/json"
	"fmt"
	"main/infrastructure"
	"main/services/common"
	"strings"

	"github.com/hashicorp/consul/api"
)

type ConsulClient struct {
	logger *common.Logger
}

func BuildConsulClient(logger *common.Logger, storageName string) *ConsulClient {
	result := &ConsulClient{
		logger: logger,
	}

	result.init(storageName)
	return result
}

func (service *ConsulClient) Get(key string) interface{} {
	pair, _, err := kvClient.Get(fullStorageName, nil)
	if err != nil {
		panic(err)
	}

	var results map[string]interface{}
	if err := json.Unmarshal(pair.Value, &results); err != nil {
		service.logger.LogError(err, err.Error())
		panic(err)
	}

	return results[key]
}

func getFullStorageName(storageName string) string {
	name := infrastructure.GetEnvironmentName()
	lowerCasedName := strings.ToLower(name)
	lowerCasedStorageName := strings.ToLower(storageName)

	return fmt.Sprintf("/%s/%s", lowerCasedName, lowerCasedStorageName)
}

func (service *ConsulClient) init(storageName string) {
	consulAddress := infrastructure.GetEnvironmentVariable("PNKL_CONSUL_ADDR")
	consulToken := infrastructure.GetEnvironmentVariable("PNKL_CONSUL_TOKEN")

	fullStorageName = getFullStorageName(storageName)

	client, err := api.NewClient(&api.Config{
		Address: consulAddress,
		Scheme:  "http",
		Token:   consulToken,
	})
	if err != nil {
		service.logger.LogError(err, "Listen error: %s\n", err)
		panic(err)
	}

	kvClient = client.KV()
}

var kvClient *api.KV
var fullStorageName string
