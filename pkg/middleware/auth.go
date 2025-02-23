package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func IsAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		jwtToken := strings.TrimPrefix(authorizationHeader, "Bearer ")

		if jwtToken != "" {
			fmt.Println("JWT Token", jwtToken)
		}

		next.ServeHTTP(w, r)
	})
}
