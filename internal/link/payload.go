package link

import "time"

type LinkCreateRequest struct {
	Url string `json:"url" validate:"required,url"`
}

type LinkUpdateRequest struct {
	Url  string `json:"url" validate:"required"`
	Hash string `json:"hash"`
}

type LinkResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Hash      string    `json:"hash"`
	Url       string    `json:"url"`
}
