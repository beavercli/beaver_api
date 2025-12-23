package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/beavercli/beaver_api/internal/service"
)

const UserContextKey = "UserContextKey"

func (s *server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ht := r.Header.Get("Authorization")
		if ht == "" {
			jsonError(w, http.StatusUnauthorized, "Request is missing an Authorization header")
			return
		}

		token := strings.Split(ht, " ")

		if token[0] != "Bearer" {
			jsonError(w, http.StatusUnauthorized, fmt.Sprintf("Token type %s is not supported", token[0]))
			return
		}

		user, err := s.service.AuthUser(r.Context(), service.AccessToken, token[1])
		if err != nil {
			jsonError(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, user)
		r = r.WithContext(ctx)

		next(w, r)
	}

}
