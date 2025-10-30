package response

import (
	"encoding/json"
	"errors"
	"net/http"

	apperror "github.com/easy-comerce/backend/pkg/app_error"
	c "github.com/easy-comerce/backend/pkg/constants"
	"github.com/easy-comerce/backend/pkg/validator"
)

type Response struct {
	Success bool    `json:"success"`
	Data    *any    `json:"data,omitempty"`
	Message *string `json:"message,omitempty"`
	Error   *Error  `json:"error,omitempty"`
}

type Error struct {
	Details []validator.ValidationError `json:"details,omitempty"`
}

func SendErrorJSON(w http.ResponseWriter, message string, statusCode ...int) {
	code := http.StatusBadRequest
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	response := Response{
		Success: false,
		Message: &message,
	}
	sendJSON(w, response, code)
}

func SendValidationErrorJSON(w http.ResponseWriter, message string, validationErrors validator.ValidationErrors, statusCode ...int) {
	code := http.StatusBadRequest
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	response := Response{
		Success: false,
		Message: &message,
		Error: &Error{
			Details: validationErrors,
		},
	}
	sendJSON(w, response, code)
}

func SendApiResponse(w http.ResponseWriter, result any, err error) {
	errorStatusCode := http.StatusBadRequest
	if err != nil && errors.As(err, new(*apperror.ServerError)) {
		errorStatusCode = http.StatusInternalServerError
	}

	SendApiResponseWithStatusCode(w, result, err, errorStatusCode)
}

func SendApiResponseWithStatusCode(w http.ResponseWriter, result any, err error, errorStatusCode int) {
	if err != nil {
		SendErrorJSON(w, err.Error(), errorStatusCode)
		return
	}

	if result != nil {
		if successMsg, ok := result.(string); ok && successMsg != "" {
			SendSuccessMsgJSON(w, successMsg)
			return
		}
	}

	SendSuccessJSON(w, result)
}

func SendResponse(w http.ResponseWriter, result any, err error, errorStatusCode int) {
	if err != nil {
		SendErrorJSON(w, err.Error(), errorStatusCode)
		return
	}

	if result != nil {
		if successMsg, ok := result.(string); ok && successMsg != "" {
			SendSuccessMsgJSON(w, successMsg)
			return
		}
	}

	SendSuccessJSON(w, result)
}

func SendSuccessJSON(w http.ResponseWriter, data any, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	response := Response{
		Success: true,
		Data:    &data,
	}
	sendJSON(w, response, code)
}

func SendSuccessMsgJSON(w http.ResponseWriter, message string, statusCode ...int) {
	code := http.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}
	response := Response{
		Success: true,
		Message: &message,
	}
	sendJSON(w, response, code)
}

func sendJSON(w http.ResponseWriter, response Response, statusCode int) {
	w.Header().Set(c.ContentType, "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
