package response

import (
	"encoding/json"
	"net/http"

	"github.com/easy-comerce/backend/pkg/models"
	"github.com/easy-comerce/backend/pkg/validator"
)

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Details []validator.ValidationError `json:"details,omitempty"`
	Message string                      `json:"message"`
	Code    string                      `json:"code,omitempty"`
}

func SendSuccessJSON(w http.ResponseWriter, data any, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	response := Response{
		Success: true,
		Data:    data,
	}
	sendJSON(w, response, code)
}

func SendErrorJSON(w http.ResponseWriter, message string, statusCode ...int) {
	code := http.StatusInternalServerError
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	response := Response{
		Success: false,
		Error: &Error{
			Message: message,
		},
	}
	sendJSON(w, response, code)
}

func SendDeleteJSON(w http.ResponseWriter, message string, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	response := Response{
		Success: true,
		Data: map[string]string{
			"message": message,
		},
	}
	sendJSON(w, response, code)
}

func SendCommonResponseJSON(w http.ResponseWriter, message string, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	response := Response{
		Success: true,
		Data: models.CommonResponse{
			Message: message,
		},
	}
	sendJSON(w, response, code)
}

func SendValidationErrorJSON(w http.ResponseWriter, message string, validationErrors []validator.ValidationError, statusCode ...int) {
	code := http.StatusBadRequest
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	response := Response{
		Success: false,
		Error: &Error{
			Message: message,
			Details: validationErrors,
		},
	}
	sendJSON(w, response, code)
}

func sendJSON(w http.ResponseWriter, response Response, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
