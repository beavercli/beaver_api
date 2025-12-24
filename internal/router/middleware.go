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

		var tokenType service.TokenType

		switch token[0] {
		case "Bearer":
			tokenType = service.AccessToken
			break
		case "Session":
			tokenType = service.SessionToken
			break
		default:
			jsonError(w, http.StatusUnauthorized, fmt.Sprintf("Token type %s is not supported", token[0]))
			return
		}

		userID, err := s.service.AuthUser(r.Context(), tokenType, token[1])
		if err != nil {
			jsonError(w, http.StatusUnauthorized, err.Error())
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, userID)
		r = r.WithContext(ctx)

		next(w, r)
	}

}
