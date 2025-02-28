package link

import (
	"fmt"
	"net/http"
	"strconv"
	config "url-shortener/configs"
	"url-shortener/internal/stat"
	"url-shortener/pkg/middleware"
	"url-shortener/pkg/req"
	"url-shortener/pkg/res"

	"gorm.io/gorm"
)

type LinkHandlerDeps struct {
	*LinkRepository
	*stat.StatRepository
	*config.Config
}

type LinkHandler struct {
	*LinkRepository
	*stat.StatRepository
}

func NewLinkHandler(router *http.ServeMux, deps LinkHandlerDeps) {
	handler := &LinkHandler{
		LinkRepository: deps.LinkRepository,
		StatRepository: deps.StatRepository,
	}

	router.Handle("GET /link", middleware.IsAuthed(handler.GetLinks(), deps.Config))
	router.Handle("POST /link", middleware.IsAuthed(handler.CreateLink(), deps.Config))
	router.HandleFunc("/{hash}", handler.GoToLink())
	router.Handle("PUT /link/{id}", middleware.IsAuthed(handler.UpdateLink(), deps.Config))
	router.Handle("DELETE /link/{id}", middleware.IsAuthed(handler.DeleteLink(), deps.Config))
}

// @Summary Получение списка коротких ссылок
// @Description Получает список коротких ссылок
// @Tags Link
// @Accept json
// @Produce json
// @Success 200 {object} res.PaginationResponse[Link]
// @Router /link [post]
func (handler *LinkHandler) GetLinks() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
		if err != nil {
			res.Error(w, ErrInvalidLimit, http.StatusBadRequest)
			return
		}

		offset, err := strconv.Atoi(r.URL.Query().Get("offset"))
		if err != nil {
			res.Error(w, ErrInvalidOffset, http.StatusBadRequest)
			return
		}

		links := handler.LinkRepository.GetLinks(limit, offset)
		count := handler.LinkRepository.GetCount()

		linksResponse := []LinkResponse{}

		for _, link := range links {
			linksResponse = append(linksResponse, LinkResponse{
				ID:        link.ID,
				Hash:      link.Hash,
				Url:       link.Url,
				CreatedAt: link.CreatedAt,
				UpdatedAt: link.UpdatedAt,
			})
		}

		res.Json(w, res.PaginationResponse[LinkResponse]{
			Items: linksResponse,
			Count: count,
		}, http.StatusOK)
	}
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
			res.Error(w, err.Error(), http.StatusBadRequest)
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
			res.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, LinkResponse{
			ID:        createdLink.ID,
			Hash:      createdLink.Hash,
			Url:       createdLink.Url,
			CreatedAt: createdLink.CreatedAt,
			UpdatedAt: createdLink.UpdatedAt,
		}, http.StatusCreated)
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
			res.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		handler.StatRepository.AddClick(link.ID)

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
		user, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		if !ok {
			panic("Context fail")
		}

		fmt.Println(user)

		payload, err := req.HandleBody[LinkUpdateRequest](w, r)
		if err != nil {
			res.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		idString := r.PathValue("id")
		id, err := strconv.ParseUint(idString, 10, 32)
		if err != nil {
			res.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		link, err := handler.LinkRepository.Update(&Link{
			Model: gorm.Model{ID: uint(id)},
			Url:   payload.Url,
			Hash:  payload.Hash,
		})
		if err != nil {
			res.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, LinkResponse{
			ID:        link.ID,
			Hash:      link.Hash,
			Url:       link.Url,
			CreatedAt: link.CreatedAt,
			UpdatedAt: link.UpdatedAt,
		}, http.StatusOK)
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
			res.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = handler.LinkRepository.FindById(uint(id))
		if err != nil {
			res.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		err = handler.LinkRepository.Delete(uint(id))
		if err != nil {
			res.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		res.Json(w, "Link deleted!", http.StatusOK)
	}
}
