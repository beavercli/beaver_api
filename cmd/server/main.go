package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beavercli/beaver_api/common/config"
	_ "github.com/beavercli/beaver_api/docs"
	"github.com/beavercli/beaver_api/internal/router"
)

// @title           Beaver API
// @version         1.0
// @description     Code snippets API

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	cfg := config.New()

	srv := router.New(router.Config{
		Addr:         cfg.Server.Addr,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println(err)
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Server stopped")
}
