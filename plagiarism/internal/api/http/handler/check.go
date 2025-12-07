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
		respondMethodNotAllowed(w, "only POST is allowed")
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
		respondValidationError(w, "failed to parse request body")
		return
	}

	if request.SubmissionID == "" {
		respondValidationError(w, "submission_id is required")
		return
	}

	if request.WorkID == "" {
		respondValidationError(w, "work_id is required")
		return
	}

	if h.useCase == nil {
		respondError(w, usecase.ErrWorkerUnavailable)
		return
	}

	resp, err := h.useCase.StartCheck(r.Context(), request.SubmissionID, request.WorkID)
	if err != nil {
		respondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(resp)
}
