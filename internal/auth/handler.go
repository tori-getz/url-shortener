package auth

import (
	"fmt"
	"net/http"
	config "url-shortener/configs"
	"url-shortener/pkg/req"
	"url-shortener/pkg/res"
)

type AuthHandlerDeps struct {
	*config.Config
}

type AuthHandler struct {
	*config.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}

	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := req.HandleBody[LoginRequest](w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}

		fmt.Println(payload)

		response := LoginResponse{
			Token: handler.Config.Auth.Secret,
		}

		res.Json(w, response, 200)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := req.HandleBody[RegisterRequest](w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}

		fmt.Println(payload)

		response := RegisterResponse{
			Name: payload.Name,
		}

		res.Json(w, response, 200)
	}
}
