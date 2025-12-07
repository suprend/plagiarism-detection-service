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
		respondMethodNotAllowed(w, "only GET method is allowed")
		return
	}

	path := r.URL.Path
	if !strings.HasPrefix(path, "/works/") || !strings.HasSuffix(path, "/reports") {
		respondValidationError(w, "invalid path format")
		return
	}

	path = strings.TrimPrefix(path, "/works/")
	path = strings.TrimSuffix(path, "/reports")

	if path == "" {
		respondValidationError(w, "work_id is required in path")
		return
	}

	workID := path

	if h.useCase == nil {
		respondError(w, usecase.ErrWorkerUnavailable)
		return
	}

	resp, err := h.useCase.GetReportsByWork(r.Context(), workID)
	if err != nil {
		if errors.Is(err, usecase.ErrCheckNotFound) {
			respondError(w, err)
			return
		}
		respondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
