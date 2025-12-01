package router

import (
	"fmt"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Config struct {
	Addr         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type server struct {
}

func New(cfg Config) *http.Server {
	mux := http.NewServeMux()
	s := &server{}

	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("GET /api/v1/snippets/random", s.handleGetRandomSnippet)
	mux.HandleFunc("POST /api/v1/snippets", s.handleCreateSnippet)

	mux.HandleFunc("GET /api/v1/tags", s.handleListTags)
	mux.HandleFunc("GET /api/v1/languages", s.handleListLanguages)

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Addr, cfg.Port),
		Handler:      mux,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}

// @Summary      Health check
// @Description  Returns health status of the API
// @Tags         health
// @Produce      plain
// @Success      200  {string}  string  "healthy"
// @Router       /health [get]
func (s *server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "healthy")
}
