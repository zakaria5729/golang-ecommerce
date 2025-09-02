package response

import (
	"encoding/json"
	"net/http"

	"github.com/easy-comerce/backend/pkg/validator"
)

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   *Error `json:"error,omitempty"`
}

type Error struct {
	Message string                      `json:"message"`
	Code    string                      `json:"code,omitempty"`
	Details []validator.ValidationError `json:"details,omitempty"`
}

func JSON(w http.ResponseWriter, data interface{}) {
	response := Response{
		Success: true,
		Data:    data,
	}
	sendJSON(w, response, http.StatusOK)
}

func JSONError(w http.ResponseWriter, message string, statusCode int) {
	response := Response{
		Success: false,
		Error: &Error{
			Message: message,
		},
	}
	sendJSON(w, response, statusCode)
}

func JSONValidationError(w http.ResponseWriter, message string, validationErrors validator.ValidationErrors) {
	response := Response{
		Success: false,
		Error: &Error{
			Message: message,
			Code:    "VALIDATION_ERROR",
			Details: validationErrors,
		},
	}
	sendJSON(w, response, http.StatusBadRequest)
}

func JSONCreated(w http.ResponseWriter, data interface{}) {
	response := Response{
		Success: true,
		Data:    data,
	}
	sendJSON(w, response, http.StatusCreated)
}

func JSONNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func sendJSON(w http.ResponseWriter, response Response, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
