package auth

import (
	"net/http"
	config "url-shortener/configs"
	"url-shortener/pkg/req"
	"url-shortener/pkg/res"
)

type AuthHandlerDeps struct {
	*config.Config
	*AuthService
}

type AuthHandler struct {
	*config.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}

	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

// @Summary Авторизация пользователя
// @Description Авторизует пользователя и возвращает токен
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body LoginRequest true "Данные для авторизации"
// @Success 200 {object} LoginResponse
// @Router /auth/login [post]
func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := req.HandleBody[LoginRequest](w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}

		response := LoginResponse{
			Token: handler.Config.Auth.Secret,
		}

		res.Json(w, response, 200)
	}
}

// @Summary Регистрация пользователя
// @Description Регистрирует пользователя и возвращает токен
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body RegisterRequest true "Данные для авторизации"
// @Success 201 {object} RegisterResponse
// @Router /auth/register [post]
func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := req.HandleBody[RegisterRequest](w, r)
		if err != nil {
			res.Json(w, err.Error(), 400)
			return
		}

		email, err := handler.AuthService.Register(payload.Name, payload.Email, payload.Password)

		res.Json(w, RegisterResponse{
			Token: email,
		}, http.StatusCreated)
	}
}
