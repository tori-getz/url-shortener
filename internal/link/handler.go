package link

import (
	"net/http"
	"strconv"
	"url-shortener/pkg/middleware"
	"url-shortener/pkg/req"
	"url-shortener/pkg/res"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	LinkRepository *LinkRepository
}

type LinkHandler struct {
	LinkRepository *LinkRepository
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
	}

	router.Handle("POST /link", middleware.IsAuthed(handler.CreateLink()))
	router.HandleFunc("/{hash}", handler.GoToLink())
	router.Handle("PUT /link/{id}", middleware.IsAuthed(handler.UpdateLink()))
	router.Handle("DELETE /link/{id}", middleware.IsAuthed(handler.DeleteLink()))
}

// @Summary Создание короткой ссылки
// @Description Создает короткую ссылку
// @Tags Link
// @Accept json
// @Produce json
// @Param input body LinkCreateRequest true "Данные для создания ссылки"
// @Success 201 {object} Link
// @Router /link [post]
func (handler *LinkHandler) CreateLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := req.HandleBody[LinkCreateRequest](w, r)
		if err != nil {
			res.Json(w, err.Error(), http.StatusBadRequest)
			return
		}

		link := NewLink(payload.Url)

		for {
			existedLink, _ := handler.LinkRepository.FindByHash(link.Hash)

			if existedLink == nil {
				break
			}

			link.GenerateHash()
		}

		createdLink, err := handler.LinkRepository.Create(link)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, createdLink, http.StatusCreated)
	}
}

// @Summary Переход по короткой ссылке
// @Description Переходит по короткой ссылке
// @Tags Link
// @Success 200
// @Router /{hash} [get]
func (handler *LinkHandler) GoToLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")

		link, err := handler.LinkRepository.FindByHash(hash)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Redirect(w, r, link.Url, http.StatusTemporaryRedirect)
	}
}

// @Summary Обновление короткой ссылки
// @Description Обновляет короткую ссылку
// @Tags Link
// @Accept json
// @Produce json
// @Param input body LinkUpdateRequest true "Данные для обновления ссылки"
// @Success 200 {object} Link
// @Router /link/{hash} [put]
func (handler *LinkHandler) UpdateLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := req.HandleBody[LinkUpdateRequest](w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   payload.Url,
			Hash:  payload.Hash,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		res.Json(w, link, 201)
	}
}

// @Summary Удаление короткой ссылки
// @Description Удаляет короткую ссылку
// @Tags Link
// @Accept json
// @Produce json
// @Success 200
// @Router /link/{hash} [delete]
func (handler *LinkHandler) DeleteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = handler.LinkRepository.FindById(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = handler.LinkRepository.Delete(uint(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, "Link deleted!", 200)
	}
}
