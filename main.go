package main

import (
	"fmt"
	"main/startup"
	"main/utils"
	"net/http"
)

func main() {
	environmentName := utils.GetEnvironmentVariable("GO_ENVIRONMENT")
	fmt.Println("Artist Updater API is running as " + environmentName)

	app := startup.Configure()

	hostAddress := utils.GetEnvironmentVariable("GIN_ADDR")
	hostPort := utils.GetEnvironmentVariable("GIN_PORT")
	http.ListenAndServe(fmt.Sprintf("%s:%s", hostAddress, hostPort), app)
}
