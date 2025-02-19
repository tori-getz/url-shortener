package middleware

import "net/http"

type MiddlewareFunc = func(next http.Handler) http.Handler

func Compose(middlewares ...MiddlewareFunc) MiddlewareFunc {
	return func(finalHandler http.Handler) http.Handler {
		handler := finalHandler

		for i := len(middlewares) - 1; i >= 0; i-- {
			handler = middlewares[i](handler)
		}

		return handler
	}
}
