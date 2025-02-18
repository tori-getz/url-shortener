package req

import "github.com/go-playground/validator/v10"

func Validate[T any](payload T) error {
	validate := validator.New()
	return validate.Struct(payload)
}
