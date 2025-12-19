package router

import (
	"net/http"

	"github.com/beavercli/beaver_api/common/config"
	"github.com/beavercli/beaver_api/internal/service"
	httpSwagger "github.com/swaggo/http-swagger"
)

type server struct {
	service *service.Service
}

func New(cfg config.Server, service *service.Service) *http.Server {
	mux := http.NewServeMux()

	s := &server{
		service: service,
	}

	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	mux.HandleFunc("GET /api/v1/snippets/{SnippetID}", s.handleGetSnippet)
	mux.HandleFunc("GET /api/v1/snippets", s.handleListSnippets)
	mux.HandleFunc("POST /api/v1/snippets", s.handleIngestSnippet)

	mux.HandleFunc("GET /api/v1/tags", s.handleListTags)
	mux.HandleFunc("GET /api/v1/languages", s.handleListLanguages)
	mux.HandleFunc("GET /api/v1/contributors", s.handleListContributors)

	mux.HandleFunc("POST /api/v1/service-access-tokens", s.handleCreateServiceAccessToken)   // todo
	mux.HandleFunc("GET /api/v1/service-access-tokens", s.handleGetServiceAccessToken)       // todo
	mux.HandleFunc("DELETE /api/v1/service-access-tokens", s.handleDeleteServiceAccessToken) // todo

	mux.HandleFunc("POST /auth/github/login", s.handleGithubLogin)              // todo
	mux.HandleFunc("POST /auth/github/device/poll", s.handleGitHubDeviceStatus) // todo
	mux.HandleFunc("POST /auth/refresh", s.handleTokenRotate)                   // todo
	mux.HandleFunc("POST /auth/logout", s.handleLogout)                         // todo
	mux.HandleFunc("GET /auth/me", s.handleMe)                                  // todo

	return &http.Server{
		Addr:         cfg.Addr,
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
	jsonResponse(w, http.StatusOK, MessageResponse{Message: "healthy"})
}
