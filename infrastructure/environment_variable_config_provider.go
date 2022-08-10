package infrastructure

import (
	"bufio"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

func init() {
	localVariables := getLocalVariables()
	systemVariables := getSystemVariables()

	totalLength := len(localVariables) + len(systemVariables)

	variables = make(map[string]string, totalLength)
	for key, value := range localVariables {
		variables[key] = value
	}

	for key, value := range systemVariables {
		variables[key] = value
	}
}

func GetEnvironmentName() string {
	name := GetEnvironmentVariable("GO_ENVIRONMENT")
	if name != "" {
		return name
	}

	return "Development"
}

func GetEnvironmentVariable(name string) string {
	if result, ok := variables[name]; ok {
		return result
	}

	log.Warn().Msgf("No configuration value has been found for '%s'", name)
	return ""
}

func GetLogLevel() string {
	level := GetEnvironmentVariable("LOG_LEVEL")
	if level != "" {
		return level
	}

	return "Error"
}

func getLocalVariables() map[string]string {
	file, err := os.Open(".env")
	if err != nil {
		log.Warn().Msgf("Can't read configuration file: %s", err)
		return make(map[string]string, 0)
	}

	fileScanner := bufio.NewScanner(file)
	fileScanner.Split(bufio.ScanLines)

	slice := make([]string, 0)
	for fileScanner.Scan() {
		slice = append(slice, fileScanner.Text())
	}

	file.Close()

	return splitKeysAndValues(slice)
}

func getSystemVariables() map[string]string {
	variables := os.Environ()
	return splitKeysAndValues(variables)
}

func splitKeysAndValues(source []string) map[string]string {
	result := make(map[string]string, len(source))

	for _, value := range source {
		split := strings.SplitN(value, "=", 2)
		result[split[0]] = split[1]
	}

	return result
}

var variables map[string]string
