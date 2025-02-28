package middleware

import (
	"context"
	"net/http"
	"strings"
	"url-shortener/configs"
	"url-shortener/pkg/jwt"
	"url-shortener/pkg/res"
)

const ErrUnauthorized = "UNAUTHORIZED"

type emailKey string

const (
	ContextEmailKey = "ContextEmailKey"
)

func IsAuthed(next http.Handler, config *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		jwtToken := strings.TrimPrefix(authorizationHeader, "Bearer ")

		isValid, payload := jwt.NewJwt(config.Auth.Secret).Parse(jwtToken)

		if !isValid {
			res.Error(w, ErrUnauthorized, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextEmailKey, payload.Email)
		reqWithUser := r.WithContext(ctx)

		next.ServeHTTP(w, reqWithUser)
	})
}
