package handler

import (
	"encoding/json"
	"net/http"

	apperr "wordcloud/internal/common/errors"
)

type apiError struct {
	status  int
	code    apperr.Code
	message string
}

func respondError(w http.ResponseWriter, err error) {
	apiErr := mapAPIError(err)
	writeJSONError(w, apiErr.status, string(apiErr.code), apiErr.message)
}

func respondValidationError(w http.ResponseWriter, message string) {
	writeJSONError(w, http.StatusBadRequest, string(apperr.CodeValidation), message)
}

func respondMethodNotAllowed(w http.ResponseWriter, message string) {
	writeJSONError(w, http.StatusMethodNotAllowed, "method_not_allowed", message)
}

func mapAPIError(err error) apiError {
	code := apperr.CodeOf(err)
	message := apperr.Message(err)
	if message == "" {
		message = defaultMessageForCode(code)
	}

	switch code {
	case apperr.CodeValidation:
		return apiError{status: http.StatusBadRequest, code: code, message: message}
	case apperr.CodeNotFound:
		return apiError{status: http.StatusNotFound, code: code, message: message}
	case apperr.CodeStorage:
		return apiError{status: http.StatusBadGateway, code: code, message: message}
	default:
		return apiError{status: http.StatusInternalServerError, code: apperr.CodeInternal, message: defaultMessageForCode(apperr.CodeInternal)}
	}
}

func defaultMessageForCode(code apperr.Code) string {
	switch code {
	case apperr.CodeValidation:
		return "validation error"
	case apperr.CodeNotFound:
		return "resource not found"
	case apperr.CodeStorage:
		return "storage error"
	default:
		return "internal error"
	}
}

func writeJSONError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error":   code,
		"message": message,
	})
}
