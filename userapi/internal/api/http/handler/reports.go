package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"userapi/internal/application/usecase"
)

type ReportsHandler struct {
	useCase *usecase.ReportsUseCase
}

func NewReportsHandler(uc *usecase.ReportsUseCase) *ReportsHandler {
	return &ReportsHandler{useCase: uc}
}

func (h *ReportsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "only GET is allowed")
		return
	}

	workID, ok := extractWorkID(r.URL.Path, "/reports")
	if !ok {
		writeError(w, http.StatusBadRequest, "invalid_path", "expected /works/{work_id}/reports")
		return
	}

	resp, err := h.useCase.GetByWork(r.Context(), workID)
	if err != nil {
		if errors.Is(err, usecase.ErrReportNotFound) {
			writeError(w, http.StatusNotFound, "not_found", "report not found")
			return
		}
		writeError(w, http.StatusBadGateway, "reports_failed", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
