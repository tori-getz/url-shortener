package res

import "net/http"

type ErrorResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func Error(w http.ResponseWriter, message string, statusCode int) {
	errorResponse := ErrorResponse{
		Message: message,
		Code:    statusCode,
	}

	Json(w, errorResponse, statusCode)
}
