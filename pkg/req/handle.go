package req

import (
	"net/http"
)

func HandleBody[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	payload, err := Decode[T](r.Body)
	if err != nil {
		return nil, err
	}

	err = Validate(payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
