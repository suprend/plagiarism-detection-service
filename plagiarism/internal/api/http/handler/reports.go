package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"plagiarism/internal/application/usecase"
)

type ReportsHandler struct {
	useCase usecase.CheckUseCase
}

func NewReportsHandler(uc usecase.CheckUseCase) *ReportsHandler {
	return &ReportsHandler{useCase: uc}
}

func (h *ReportsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET method is allowed")
		return
	}

	path := r.URL.Path
	if !strings.HasPrefix(path, "/works/") || !strings.HasSuffix(path, "/reports") {
		h.writeError(w, http.StatusBadRequest, "validation_error", "invalid path format")
		return
	}

	path = strings.TrimPrefix(path, "/works/")
	path = strings.TrimSuffix(path, "/reports")

	if path == "" {
		h.writeError(w, http.StatusBadRequest, "validation_error", "work_id is required in path")
		return
	}

	workID := path

	if h.useCase == nil {
		h.writeError(w, http.StatusInternalServerError, "not_configured", "check use case is not configured")
		return
	}

	resp, err := h.useCase.GetReportsByWork(r.Context(), workID)
	if err != nil {
		if errors.Is(err, usecase.ErrCheckNotFound) {
			h.writeError(w, http.StatusNotFound, "not_found", "report not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "get_reports_failed", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func (h *ReportsHandler) writeError(w http.ResponseWriter, statusCode int, errorType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]string{
		"error":   errorType,
		"message": message,
	})
}
