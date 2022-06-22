package utils

import (
	"github.com/spf13/viper"
)

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}

func GetEnvironmentVariable(name string) string {
	if !viper.IsSet(name) {
		return ""
	}

	return viper.GetString(name)
}
