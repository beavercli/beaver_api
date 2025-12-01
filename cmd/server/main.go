package main

import (
	"fmt"
	"time"

	_ "github.com/beavercli/beaver_api/docs"
	"github.com/beavercli/beaver_api/internal/router"
)

// @title           Beaver API
// @version         1.0
// @description     Code snippets API

// @host      localhost:8080
// @BasePath  /api/v1
func main() {
	srv := router.New(router.Config{
		Addr:         "127.0.0.1",
		Port:         8080,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	})

	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
