package link

import (
	"fmt"
	"net/http"
)

type LinkHandler struct{}

func NewLinkHandler(router *http.ServeMux) {
	handler := &LinkHandler{}

	router.HandleFunc("POST /link", handler.CreateLink())
	router.HandleFunc("/{hash}", handler.ReadLink())
	router.HandleFunc("PUT /link/{id}", handler.UpdateLink())
	router.HandleFunc("DELETE /link/{id}", handler.DeleteLink())
}

func (handler *LinkHandler) CreateLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (handler *LinkHandler) ReadLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		hash := r.PathValue("hash")
		fmt.Println(hash)
	}
}

func (handler *LinkHandler) UpdateLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func (handler *LinkHandler) DeleteLink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
