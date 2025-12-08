package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beavercli/beaver_api/common/config"
	"github.com/beavercli/beaver_api/common/database"
	_ "github.com/beavercli/beaver_api/docs"
	"github.com/beavercli/beaver_api/internal/router"
	"github.com/beavercli/beaver_api/internal/service"
	"github.com/beavercli/beaver_api/internal/storage"
)

// @title           Beaver API
// @version         1.0
// @description     Code snippets API
// @host      localhost:8080
func main() {
	ctx := context.Background()
	cfg := config.New()

	pool, err := database.New(ctx, cfg.DB)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	storage := storage.New(pool)
	service := service.New(storage)
	server := router.New(cfg.Server, service)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println(err)
			panic(err)
		}
	}()

	fmt.Println("Starting API server")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
	fmt.Println("Server stopped")
}
