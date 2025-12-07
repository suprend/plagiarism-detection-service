package handler

import (
	"net/http"

	"wordcloud/internal/application/usecase"
)

type WordcloudHandler struct {
	useCase *usecase.WordcloudService
}

func NewWordcloudHandler(uc *usecase.WordcloudService) *WordcloudHandler {
	return &WordcloudHandler{useCase: uc}
}

func (h *WordcloudHandler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respondMethodNotAllowed(w, "only GET is allowed")
		return
	}

	submissionID := r.URL.Query().Get("submission_id")
	if submissionID == "" {
		respondValidationError(w, "submission_id is required")
		return
	}

	img, err := h.useCase.Generate(r.Context(), submissionID)
	if err != nil {
		respondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(img)
}
