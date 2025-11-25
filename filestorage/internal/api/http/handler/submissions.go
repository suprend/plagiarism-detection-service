package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"filestorage/internal/application/usecase"
)

type SubmissionsHandler struct {
	getSubmissionsUseCase *usecase.GetSubmissionsUseCase
}

func NewSubmissionsHandler(getSubmissionsUseCase *usecase.GetSubmissionsUseCase) *SubmissionsHandler {
	return &SubmissionsHandler{
		getSubmissionsUseCase: getSubmissionsUseCase,
	}
}

func (h *SubmissionsHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondMethodNotAllowed(w, "only GET method is allowed")
		return
	}

	// Получение assignment_id из query параметров
	assignmentID := r.URL.Query().Get("assignment_id")
	if assignmentID == "" {
		respondValidationError(w, "assignment_id query parameter is required")
		return
	}

	// Получаем список сдач
	submissions, err := h.getSubmissionsUseCase.GetByAssignmentID(r.Context(), assignmentID)
	if err != nil {
		log.Printf("submissions: assignment_id=%s failed: %v", assignmentID, err)
		respondError(w, err)
		return
	}

	// Формируем ответ
	submissionsResponse := make([]map[string]interface{}, 0, len(submissions))
	for _, sub := range submissions {
		submissionsResponse = append(submissionsResponse, map[string]interface{}{
			"submission_id": sub.SubmissionID.String(),
			"assignment_id": sub.AssignmentID,
			"author_id":     sub.AuthorID,
			"created_at":    sub.CreatedAt,
		})
	}

	response := map[string]interface{}{
		"submissions": submissionsResponse,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
