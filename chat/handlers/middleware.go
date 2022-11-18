package handlers

import (
	"context"
	"github.com/sornick01/http_chat/chat"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	useCase chat.UseCase
}

func NewAuthMiddleware(useCase chat.UseCase) *AuthMiddleware {
	return &AuthMiddleware{
		useCase: useCase,
	}
}

func (am *AuthMiddleware) Auth(handler http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if headerParts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		user, err := am.useCase.ParseToken(r.Context(), headerParts[1])
		if err != nil {
			if err == chat.ErrInvalidAccessToken {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		userCtx := context.WithValue(r.Context(), chat.CtxUserKey, user)
		handler.ServeHTTP(w, r.WithContext(userCtx))
	}

	return http.HandlerFunc(fn)
}
