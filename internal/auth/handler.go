package auth

import (
	"net/http"
	config "url-shortener/configs"
	"url-shortener/pkg/jwt"
	"url-shortener/pkg/req"
	"url-shortener/pkg/res"
)

type AuthHandlerDeps struct {
	*config.Config
	*AuthService
	*jwt.Jwt
}

type AuthHandler struct {
	*config.Config
	*AuthService
	*jwt.Jwt
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
		Jwt:         deps.Jwt,
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
		payload, err := req.HandleBody[LoginRequest](w, r)
		if err != nil {
			res.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		email, err := handler.AuthService.Login(payload.Email, payload.Password)
		if err != nil {
			res.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := handler.Jwt.Create(jwt.JwtPayload{
			Email: email,
		})
		if err != nil {
			res.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response := LoginResponse{
			Token: token,
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
			res.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		email, err := handler.AuthService.Register(payload.Name, payload.Email, payload.Password)
		if err != nil {
			res.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		token, err := handler.Jwt.Create(jwt.JwtPayload{
			Email: email,
		})
		if err != nil {
			res.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Json(w, RegisterResponse{
			Token: token,
		}, http.StatusCreated)
	}
}
