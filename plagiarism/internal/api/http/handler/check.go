package handler

import (
	"encoding/json"
	"net/http"

	"plagiarism/internal/application/usecase"
)

type CheckHandler struct {
	useCase usecase.CheckUseCase
}

func NewCheckHandler(uc usecase.CheckUseCase) *CheckHandler {
	return &CheckHandler{useCase: uc}
}

func (h *CheckHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only POST is allowed")
		return
	}
	h.handleStart(w, r)
}

func (h *CheckHandler) handleStart(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SubmissionID string `json:"submission_id"`
		WorkID       string `json:"work_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.writeError(w, http.StatusBadRequest, "parse_error", "failed to parse request body")
		return
	}

	if request.SubmissionID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "submission_id is required")
		return
	}

	if request.WorkID == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "work_id is required")
		return
	}

	if h.useCase == nil {
		h.writeError(w, http.StatusInternalServerError, "not_configured", "check use case is not configured")
		return
	}

	resp, err := h.useCase.StartCheck(r.Context(), request.SubmissionID, request.WorkID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "start_check_failed", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}

func (h *CheckHandler) writeError(w http.ResponseWriter, statusCode int, errorType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   errorType,
		"message": message,
	})
}
