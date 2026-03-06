package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"perftest-target/config"
	"perftest-target/router"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

const appVersion = "2026-03-06"

func main() {

	printBanner()

	// read configuration
	config.ReadConfiguration()

	// initialize GIN router and register our routes
	router := router.New()

	// start the server.
	bind := config.GetEnvString("HTTP_BIND_ADDR")
	if bind == "" {
		bind = "0.0.0.0:3000"
		logrus.Warnf("HTTP_BIND_ADDR not set, defaulting to %s", bind)
	}

	srv := &http.Server{
		Addr:    bind,
		Handler: router,
	}

	// start server in a goroutine
	go func() {
		logrus.Infof("Server listening on http://%s", bind)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("There was an error with the http server: %v", err)
		}
	}()

	// wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}

	logrus.Info("Server exited")
}

func printBanner() {
	banner := `
██████╗ ███████╗██████╗ ███████╗ ████████╗███████╗███████╗████████╗
██╔══██╗██╔════╝██╔══██╗██╔════╝ ╚══██╔══╝██╔════╝██╔════╝╚══██╔══╝
██████╔╝█████╗  ██████╔╝█████╗█████╗██║   █████╗  ███████╗   ██║   
██╔═══╝ ██╔══╝  ██╔══██╗██╔══╝╚════╝██║   ██╔══╝  ╚════██║   ██║   
██║     ███████╗██║  ██║██║         ██║   ███████╗███████║   ██║   
╚═╝     ╚══════╝╚═╝  ╚═╝╚═╝         ╚═╝   ╚══════╝╚══════╝   ╚═╝   
                                                                   
App version: %s
`

	b := fmt.Sprintf(banner, appVersion)
	fmt.Println(b)
}
