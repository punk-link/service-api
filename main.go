package main

import (
	"context"
	"errors"
	"fmt"
	"main/startup"
	"main/utils"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	environmentName := utils.GetEnvironmentVariable("GO_ENVIRONMENT")
	fmt.Printf("Artist Updater API is running as '%s'\n", environmentName)

	hostAddress := utils.GetEnvironmentVariable("GIN_ADDR")
	hostPort := utils.GetEnvironmentVariable("GIN_PORT")

	app := startup.Configure()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", hostAddress, hostPort),
		Handler: app,
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Listen error: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Sutting down...")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFunc()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %s", err)
	}

	fmt.Println("Exiting")
}
