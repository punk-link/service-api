package consul

import (
	"encoding/json"
	"fmt"
	"main/infrastructure"
	"strings"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/punk-link/logger"
)

type ConsulClient struct {
	logger *logger.Logger
}

func New(logger *logger.Logger, storageName string) *ConsulClient {
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

func (service *ConsulClient) GetOrSet(key string, period time.Duration) interface{} {
	now := time.Now().UTC()
	if container, ok := localStorage[key]; ok {
		if now.Before(container.Expired) {
			return container.Value
		}
	}

	if period == 0 {
		period = time.Duration(time.Minute * 5)
	}

	value := service.Get(key)
	localStorage[key] = LocalCacheContainer{
		Expired: now.Add(period),
		Value:   value,
	}

	return value
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

var fullStorageName string
var kvClient *api.KV
var localStorage map[string]LocalCacheContainer = make(map[string]LocalCacheContainer)
