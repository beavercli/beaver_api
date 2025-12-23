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
	"github.com/beavercli/beaver_api/internal/integrations/github"
	"github.com/beavercli/beaver_api/internal/router"
	"github.com/beavercli/beaver_api/internal/service"
)

// @title           Beaver API
// @version         1.0
// @description     Code snippets API
// @securitydefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host      localhost:8080
func main() {
	ctx := context.Background()
	cfg := config.New()

	pool, err := database.New(ctx, cfg.DB)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	ghCLient := github.New(cfg.OAuth.ClientID, 2*time.Second)
	service := service.New(pool, service.Config{Secret: cfg.OAuth.Secret}, ghCLient)
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
