package handler

import (
	"encoding/json"
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
		respondMethodNotAllowed(w, "only GET is allowed")
		return
	}

	workID, ok := extractWorkID(r.URL.Path, "/reports")
	if !ok {
		respondValidationError(w, "expected /works/{work_id}/reports")
		return
	}

	resp, err := h.useCase.GetByWork(r.Context(), workID)
	if err != nil {
		respondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
