package main

import (
	"fmt"

	"github.com/beavercli/beaver_api/internal/router"
)

func main() {
	srv := router.New()
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
		panic(err)
	}
}
