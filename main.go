package main

import (
	"fmt"
	"main/startup"
	"main/utils"
)

func main() {
	environmentName := utils.GetEnvironmentVariable("GO_ENVIRONMENT")
	fmt.Println("Artist Updater API is running as " + environmentName)

	app := startup.Configure()
	app.Run()
}
