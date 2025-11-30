package router

import (
	"fmt"
	"net/http"
	"time"
)

type server struct {
}

func New() *http.Server {
	mux := http.NewServeMux()
	s := &server{}

	mux.HandleFunc("/health", s.healthCheck)

	return &http.Server{
		Addr:         "127.0.0.1:8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func (s *server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "healthy")
}
